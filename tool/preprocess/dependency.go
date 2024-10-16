// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package preprocess

import (
	"bufio"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"golang.org/x/mod/modfile"
)

const (
	OtelSetupInst    = "otel_setup_inst.go"
	OtelSetupSDK     = "otel_setup_sdk.go"
	OtelRules        = "otel_rules"
	OtelUser         = "otel_user"
	OtelRuleCache    = "rule_cache"
	OtelBackups      = "backups"
	OtelBackupSuffix = ".bk"
	FuncMain         = "main"
	FuncInit         = "init"
	StdRulesPrefix   = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/"
	StdRulesPath     = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules"
	apiImport        = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

// @@ Change should sync with trampoline template
const (
	OtelGetStackDef          = "OtelGetStackImpl"
	OtelGetStackImportPath   = "runtime/debug"
	OtelGetStackAliasPkg     = "otel_runtime_debug"
	OtelGetStackImplCode     = OtelGetStackAliasPkg + ".Stack"
	OtelPrintStackDef        = "OtelPrintStackImpl"
	OtelPrintStackImportPath = "log"
	OtelPrintStackPkgAlias   = "otel_log"
	OtelPrintStackImplCode   = "func(bt []byte){ otel_log.Printf(string(bt)) }"
)

type DepProcessor struct {
	bundles          []*resource.RuleBundle // All dependent rule bundles
	sigc             chan os.Signal         // Graceful shutdown
	backups          map[string]string
	localImportPath  string
	importCandidates []string
	rule2Dir         map[*resource.InstFuncRule]string
	ruleCache        embed.FS
}

func newDepProcessor() *DepProcessor {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	return &DepProcessor{
		bundles:          []*resource.RuleBundle{},
		sigc:             sigc,
		backups:          map[string]string{},
		localImportPath:  "",
		importCandidates: nil,
		rule2Dir:         map[*resource.InstFuncRule]string{},
		ruleCache:        pkg.ExportRuleCache(),
	}
}

func (dp *DepProcessor) postProcess() {
	shared.GuaranteeInPreprocess()
	// Clean build cache as we may instrument some std packages(e.g. runtime)
	// TODO: fine-grained cache cleanup
	// err := util.RunCmd("go", "clean", "-cache")
	// if err != nil {
	// 	log.Fatalf("failed to clean cache: %v", err)
	// }

	// Using -debug? Leave all changes for debugging
	if shared.Debug {
		return
	}

	// rm -rf otel_rules
	_ = os.RemoveAll(OtelRules)

	// Restore everything we have modified during instrumentation
	err := dp.restoreBackupFiles()
	if err != nil {
		log.Fatalf("failed to restore: %v", err)
	}
}

func (dp *DepProcessor) catchSignal() {
	util.Assert(dp.sigc != nil, "sanity check")
	go func() {
		s := <-dp.sigc
		switch s {
		case syscall.SIGTERM, syscall.SIGINT:
			log.Printf("Interrupted instrumentation, cleaning up")
			dp.postProcess()
		default:
		}
	}()
}

func (dp *DepProcessor) backupFile(origin string) error {
	shared.GuaranteeInPreprocess()
	backup := filepath.Base(origin) + OtelBackupSuffix
	backup = shared.GetLogPath(filepath.Join(OtelBackups, backup))
	err := os.MkdirAll(filepath.Dir(backup), 0777)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if _, exist := dp.backups[origin]; !exist {
		err = util.CopyFile(origin, backup)
		if err != nil {
			return fmt.Errorf("failed to backup file %v: %w", origin, err)
		}
		dp.backups[origin] = backup
		log.Printf("Backup %v\n", origin)
	} else if shared.Verbose {
		log.Printf("Backup %v already exists\n", origin)
	}
	return nil
}

func (dp *DepProcessor) restoreBackupFiles() error {
	shared.GuaranteeInPreprocess()
	for origin, backup := range dp.backups {
		err := util.CopyFile(backup, origin)
		if err != nil {
			return err
		}
		log.Printf("Restore %v\n", origin)
	}
	return nil
}

func getCompileCommands() ([]string, error) {
	err := runDryBuild()
	if err != nil {
		// Tell us more about what happened in the dry run
		errLog, _ := util.ReadFile(shared.GetLogPath(DryRunLog))
		return nil, fmt.Errorf("failed to run dry build: %w\n%v", err, errLog)
	}
	dryRunLog, err := os.Open(shared.GetLogPath(DryRunLog))
	if err != nil {
		return nil, err
	}
	defer func(dryRunLog *os.File) {
		err := dryRunLog.Close()
		if err != nil {
			log.Printf("Failed to close dry run log file: %v", err)
		}
	}(dryRunLog)

	// Filter compile commands from dry run log
	compileCmds := make([]string, 0)
	scanner := bufio.NewScanner(dryRunLog)
	// 10MB should be enough to accommodate most long line
	buffer := make([]byte, 0, 10*1024*1024)
	scanner.Buffer(buffer, cap(buffer))
	for scanner.Scan() {
		line := scanner.Text()
		if shared.IsCompileCommand(line) {
			line = strings.Trim(line, " ")
			compileCmds = append(compileCmds, line)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return compileCmds, nil
}

// assembleInitCandidate assembles the candidate files that we may add init
// function to. The candidate files are the ones that have main or init
// function defined.
func (dp *DepProcessor) getImportCandidates() ([]string, error) {
	if dp.importCandidates != nil {
		return dp.importCandidates, nil
	}
	candidates := make([]string, 0)
	found := false

	// Find from build arguments e.g. go build test.go or go build cmd/app
	for _, buildArg := range shared.BuildArgs {
		// FIXME: Should we check file permission here? As we are likely to read
		// it later, which would cause fatal error if permission is not granted.

		// It's a golang file, good candidate
		if shared.IsGoFile(buildArg) {
			candidates = append(candidates, buildArg)
			found = true
			continue
		}
		// It's likely a flag, skip it
		if strings.HasPrefix("-", buildArg) {
			continue
		}

		// It's a directory, find all go files in it
		if exist, _ := util.PathExists(buildArg); exist {
			p2, err := util.ListFilesFlat(buildArg)
			if err != nil {
				// Error is tolerated here, as buildArg may be a file
				continue
			}
			for _, file := range p2 {
				if shared.IsGoFile(file) {
					candidates = append(candidates, file)
					found = true
				}
			}
		}
	}

	// Find candidates from current directory if no build arguments are provided
	if !found {
		files, err := util.ListFilesFlat(".")
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %w", err)
		}
		candidates = append(candidates, files...)
	}
	if len(candidates) > 0 {
		dp.importCandidates = candidates
	}
	return candidates, nil
}

func (dp *DepProcessor) addExplicitImport(importPaths ...string) (err error) {
	// Find out where we should forcely import our init func
	candidate, err := dp.getImportCandidates()
	if err != nil {
		return fmt.Errorf("failed to get import candidates: %w", err)
	}

	addImport := false
	for _, file := range candidate {
		if !shared.IsGoFile(file) {
			continue
		}
		astRoot, err := shared.ParseAstFromFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse ast from source: %w", err)
		}

		foundInit := shared.FindFuncDecl(astRoot, FuncInit) != nil
		if !foundInit {
			foundMain := shared.FindFuncDecl(astRoot, FuncMain) != nil
			if !foundMain {
				continue
			}
		}

		// Prepend import path to the file
		for _, importPath := range importPaths {
			shared.AddImportForcely(astRoot, importPath)
			log.Printf("Add %s import to %v", importPath, file)
		}
		addImport = true

		err = dp.backupFile(file)
		if err != nil {
			return fmt.Errorf("failed to backup file %v: %w", file, err)
		}
		_, err = shared.WriteAstToFile(astRoot, filepath.Join(file))
		if err != nil {
			return fmt.Errorf("failed to write ast to %v: %w", file, err)
		}
	}
	if !addImport {
		return fmt.Errorf("failed to add rule import, candidates are %v",
			dp.importCandidates)
	}
	return nil
}

// getModuleName returns the module name of the project by parsing go.mod file.
func getModuleName(gomod string) (string, error) {
	data, err := util.ReadFile(gomod)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	modFile, err := modfile.Parse(shared.GoModFile, []byte(data), nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse go.mod: %w", err)
	}

	moduleName := modFile.Module.Mod.Path
	return moduleName, nil
}

func (dp *DepProcessor) findLocalImportPath() error {
	// Get absolute path of current working directory
	workingDir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	// Get absolute path of go.mod directory
	gomod, err := shared.GetGoModPath()
	if err != nil {
		return fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	projectDir := filepath.Dir(gomod)
	// Replace go.mod directory with module name
	moduleName, err := getModuleName(gomod)
	if err != nil {
		return fmt.Errorf("failed to get module name: %w", err)
	}
	dp.localImportPath = strings.Replace(workingDir, projectDir, moduleName, 1)
	if shared.Verbose {
		log.Printf("Find local import path: %v", dp.localImportPath)
	}
	return nil
}

func (dp *DepProcessor) getImportPathOf(dirName string) (string, error) {
	util.Assert(dirName != "", "dirName is empty")
	if dp.localImportPath == "" {
		err := dp.findLocalImportPath()
		if err != nil {
			return "", fmt.Errorf("failed to find local import path: %w", err)
		}
	}
	return dp.localImportPath + "/" + dirName, nil
}

func (dp *DepProcessor) addOtelImports() error {
	deps := []string{}
	for _, dep := range fixedDeps {
		if dep.addImport {
			deps = append(deps, dep.dep)
		}
	}
	err := dp.addExplicitImport(deps...)
	if err != nil {
		return fmt.Errorf("failed to add otel import: %w", err)
	}
	return nil
}

func (dp *DepProcessor) preclean() {
	// err is tolerated here as this is a best-effort operation
	// Clean obsolete imports from last run
	candidate, _ := dp.getImportCandidates()
	ruleImport, _ := dp.getImportPathOf(OtelRules)
	for _, file := range candidate {
		if !shared.IsGoFile(file) {
			continue
		}
		astRoot, _ := shared.ParseAstFromFile(file)
		if astRoot == nil {
			continue
		}
		if shared.RemoveImport(astRoot, ruleImport) {
			if shared.Verbose {
				log.Printf("Remove obsolete import %v from %v",
					ruleImport, file)
			}
		}
		for _, dep := range fixedDeps {
			if !dep.addImport {
				continue
			}
			if shared.RemoveImport(astRoot, dep.dep) {
				if shared.Verbose {
					log.Printf("Remove obsolete import %v from %v",
						dep, file)
				}
			}
		}
		_, err := shared.WriteAstToFile(astRoot, file)
		if err != nil {
			log.Printf("Failed to write ast to %v: %v", file, err)
		}
	}
	// Clean otel_rules directory
	if exist, _ := util.PathExists(OtelRules); exist {
		_ = os.RemoveAll(OtelRules)
	}
}

func (dp *DepProcessor) storeRuleBundles() error {
	err := resource.StoreRuleBundles(dp.bundles)
	if err != nil {
		return fmt.Errorf("failed to store rule bundles: %w", err)
	}
	// No longer valid from now on
	dp.bundles = nil
	return nil
}

func (dp *DepProcessor) setupDeps() error {
	// Pre-clean before processing in case of any obsolete materials left
	dp.preclean()

	err := dp.addOtelImports()
	if err != nil {
		return fmt.Errorf("failed to add otel imports: %w", err)
	}

	// Pinning otel version in go.mod
	err = dp.pinDepVersion()
	if err != nil {
		return fmt.Errorf("failed to update otel: %w", err)
	}

	// Run go mod tidy first to fetch all dependencies
	err = runModTidy()
	if err != nil {
		return fmt.Errorf("failed to run mod tidy: %w", err)
	}

	// Find compile commands from dry run log
	compileCmds, err := getCompileCommands()
	if err != nil {
		return fmt.Errorf("failed to get compile commands: %w", err)
	}

	// Find used rules according to compile commands
	err = dp.matchRules(compileCmds)
	if err != nil {
		return fmt.Errorf("failed to find dependencies: %w", err)
	}

	err = dp.fetchRules()
	if err != nil {
		return fmt.Errorf("failed to fetch rules: %w", err)
	}

	// Setup rules according to compile commands
	err = dp.setupRules()
	if err != nil {
		return fmt.Errorf("failed to setup dependencies: %w", err)
	}

	// Save matched rules into file, from this point on, we no longer modify
	// the rules
	err = dp.storeRuleBundles()
	if err != nil {
		return fmt.Errorf("failed to store rule bundles: %w", err)
	}
	return nil
}
