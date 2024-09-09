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
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"

	"golang.org/x/mod/modfile"
)

const (
	OtelSetupInst          = "otel_setup_inst.go"
	OtelSetupSDK           = "otel_setup_sdk.go"
	OtelRules              = "otel_rules"
	OtelBackups            = "backups"
	OtelBackupSuffix       = ".bk"
	OtelImportPath         = "go.opentelemetry.io/otel"
	OtelBaggageImportPath  = "go.opentelemetry.io/otel/baggage"
	OtelSdkTraceImportPath = "go.opentelemetry.io/otel/sdk/trace"
)

// @@ Change should sync with trampoline template
const (
	OtelGetStackDef        = "OtelGetStackImpl"
	OtelGetStackImportPath = "runtime/debug"
	OtelGetStackAliasPkg   = "otel_runtime_debug"
	OtelGetStackImplCode   = OtelGetStackAliasPkg + ".Stack"

	OtelPrintStackDef        = "OtelPrintStackImpl"
	OtelPrintStackImportPath = "log"
	OtelPrintStackPkgAlias   = "otel_log"
	OtelPrintStackImplCode   = "func(bt []byte){ otel_log.Printf(string(bt)) }"
)

const (
	FuncMain = "main"
	FuncInit = "init"
)

type DepProcessor struct {
	bundles          []*resource.RuleBundle // All dependent rule bundles
	funcRules        []uint64               // Function should be processed separately
	generatedDeps    []string
	sigc             chan os.Signal // Graceful shutdown
	backups          map[string]string
	localImportPath  string
	importCandidates []string
}

func newDepProcessor() *DepProcessor {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	return &DepProcessor{
		bundles:          []*resource.RuleBundle{},
		funcRules:        []uint64{},
		generatedDeps:    []string{},
		sigc:             sigc,
		backups:          map[string]string{},
		localImportPath:  "",
		importCandidates: nil,
	}
}

func (dp *DepProcessor) postProcess() {
	shared.GuaranteeInPreprocess()
	// Clean build cache as we may instrument some std packages(e.g. runtime)
	// TODO: fine-grained cache cleanup
	err := util.RunCmd("go", "clean", "-cache")
	if err != nil {
		log.Fatalf("failed to clean cache: %v", err)
	}

	// Using -debug? Leave all changes for debugging
	if shared.Debug {
		return
	}

	// rm -rf otel_rules
	_ = os.RemoveAll(OtelRules)

	// Restore everything we have modified during instrumentation
	err = dp.restoreBackupFiles()
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

func (dp *DepProcessor) addGeneratedDep(dep string) {
	dp.generatedDeps = append(dp.generatedDeps, dep)
}

func getCompileCommands() ([]string, error) {
	// Befor generating compile commands, let's run go mod tidy first
	// to fetch all dependencies
	err := runModTidy()
	if err != nil {
		return nil, fmt.Errorf("failed to run mod tidy: %w", err)
	}
	err = runDryBuild()
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

func readImportPath(cmd []string) string {
	var pkg string
	for i, v := range cmd {
		if v == "-p" {
			return cmd[i+1]
		}
	}
	return pkg
}

func (dp *DepProcessor) matchRules(compileCmds []string) error {
	matcher := newRuleMatcher()
	// Find used instrumentation rule according to compile commands
	for _, cmd := range compileCmds {
		cmdArgs := strings.Split(cmd, " ")
		importPath := readImportPath(cmdArgs)
		if importPath == "" {
			return fmt.Errorf("failed to find import path: %v", cmd)
		}
		if shared.Verbose {
			log.Printf("Try to match rules for %v with %v\n",
				importPath, cmdArgs)
		}
		bundle := matcher.matchRuleBundle(importPath, cmdArgs)
		if bundle.IsValid() {
			dp.bundles = append(dp.bundles, bundle)
		} else if shared.Verbose {
			log.Printf("No match for %v", importPath)
		}
	}
	// In rare case, we might instrument functions that are not in the project
	// but introduced by InstFileRule/InstFuncRule. For instance, if InstFileRule
	// adds a foo.go file containing the Foo function, we would want to further
	// instrument that one. In such cases, we need to match rules for them again.
	for _, bundle := range dp.bundles {
		// FIXME: Support further instrumenting onEnter/onExit hook
		if len(bundle.FileRules) == 0 {
			continue
		}
		candidates := make([]string, 0)
		for _, ruleHash := range bundle.FileRules {
			rule := resource.FindFileRuleByHash(ruleHash)
			// @@ File rules that intended to REPLACE the original file should
			// not be considered as candidates because logically they do not
			// introduce any new dependencies, but APPEND mode does.
			if !rule.Replace {
				candidates = append(candidates, rule.FileName)
			}
		}
		log.Printf("Try to match additional %v for %v\n",
			candidates, bundle.ImportPath)
		newBundle := matcher.matchRuleBundle(bundle.ImportPath, candidates)
		// One rule bundle represents one import path, so we should merge
		// them together instead of adding a brand new one
		_, err := bundle.Merge(newBundle)
		if err != nil {
			return fmt.Errorf("failed to merge rule bundle: %w", err)
		}
	}
	// Save used rules to file, so that we can restore them in instrument phase
	// rather than re-matching them again
	err := resource.StoreRuleBundles(dp.bundles)
	if err != nil {
		return fmt.Errorf("failed to persist used rules: %w", err)
	}
	return nil
}

func (dp *DepProcessor) copyRule(path string, target string) error {
	text, err := resource.ReadRuleFile(path)
	if err != nil {
		return fmt.Errorf("failed to read rule file %v: %w", path, err)
	}
	if !shared.HasGoBuildComment(text) {
		log.Printf("Warning: %v does not contain //go:build ignore tag", path)
	}
	text = shared.RemoveGoBuildComment(text)
	astRoot, err := shared.ParseAstFromSource(text)
	if err != nil {
		return fmt.Errorf("failed to parse ast from source: %w", err)
	}
	// If go.mod is present, we should put all rules into otel_rules directory
	astRoot.Name.Name = OtelRules

	// Copy used rule into project
	_, err = shared.WriteAstToFile(astRoot, target)
	if err != nil {
		return fmt.Errorf("failed to write ast to %v: %w", target, err)
	}
	if shared.Verbose {
		log.Printf("Copy dependency %v to %v", path, target)
	}
	return nil
}

func isValidFilePath(path string) bool {
	// a valid file path looks like: /path/to/file.go
	re := regexp.MustCompile(`^.*\.go$`)
	return re.MatchString(path)
}

func (dp *DepProcessor) copyRules(targetDir string) (err error) {
	if len(dp.bundles) == 0 {
		return nil
	}
	// Find out which resource files we should add to project
	uniqueResources := make(map[string]string)
	for _, bundle := range dp.bundles {
		for _, funcRules := range bundle.File2FuncRules {
			// Copy resource file into project as otel_rule_\d.go
			for _, rs := range funcRules {
				dp.funcRules = append(dp.funcRules, rs...)
				for _, ruleHash := range rs {
					rule := resource.FindFuncRuleByHash(ruleHash)

					// Find files where hooks relies on
					for _, dep := range rule.FileDeps {
						util.Assert(isValidFilePath(dep), "sanity check")
						res, err := resource.FindRuleFile(dep)
						if err != nil {
							return fmt.Errorf("cannot find dep %v: %w", dep, err)
						}
						uniqueResources[res] = bundle.PackageName
					}
					// If rule inserts raw code directly, skip adding any
					// further dependencies
					if rule.UseRaw {
						continue
					}
					// Find files where hooks defines in
					resources, err := resource.FindRuleFiles(rule)
					if err != nil {
						return err
					}
					if resources == nil {
						return fmt.Errorf("can not find resource for %v", rule)
					}
					for _, res := range resources {
						uniqueResources[res] = bundle.PackageName
					}
				}
			}
		}
	}

	for path, pkgName := range uniqueResources {
		name := fmt.Sprintf("otel_rule_%s%s.go", pkgName, util.RandomString(5))
		ruleFile := filepath.Join(targetDir, name)
		err = dp.copyRule(path, ruleFile)
		if err != nil {
			return fmt.Errorf("failed to copy rule %v: %w", path, err)
		}
		dp.addGeneratedDep(ruleFile)
		shared.SaveDebugFile("", ruleFile)
	}
	return nil
}

func (dp *DepProcessor) initializeRules(pkgName, target string) (err error) {
	c := fmt.Sprintf("package %s\n", pkgName)
	imports := make(map[string]string)

	assigns := make([]string, 0)
	for _, bundle := range dp.bundles {
		if len(bundle.File2FuncRules) == 0 {
			continue
		}
		addedImport := false
		for _, funcRules := range bundle.File2FuncRules {
			for _, rs := range funcRules {
				for _, ruleHash := range rs {
					rule := resource.FindFuncRuleByHash(ruleHash)

					util.Assert(rule.OnEnter != "" || rule.OnExit != "",
						"sanity check")
					if rule.UseRaw {
						continue
					}
					var aliasPkg string
					if !addedImport {
						aliasPkg = bundle.PackageName + util.RandomString(3)
						imports[bundle.ImportPath] = aliasPkg
						addedImport = true
					} else {
						aliasPkg = imports[bundle.ImportPath]
					}
					if rule.OnEnter != "" {
						assigns = append(assigns,
							fmt.Sprintf("\t%s.%s = %s\n",
								aliasPkg,
								shared.GetVarNameOfFunc(rule.OnEnter),
								rule.OnEnter,
							),
						)
					}
					if rule.OnExit != "" {
						assigns = append(assigns,
							fmt.Sprintf(
								"\t%s.%s = %s\n",
								aliasPkg,
								shared.GetVarNameOfFunc(rule.OnExit),
								rule.OnExit,
							),
						)
					}
					assigns = append(assigns, fmt.Sprintf(
						"\t%s.%s = %s\n",
						aliasPkg,
						OtelGetStackDef,
						OtelGetStackImplCode,
					))
					assigns = append(assigns, fmt.Sprintf(
						"\t%s.%s = %s\n",
						aliasPkg,
						OtelPrintStackDef,
						OtelPrintStackImplCode,
					))
				}
			}
		}
	}

	// Imports
	if len(assigns) > 0 {
		imports[OtelPrintStackImportPath] = OtelPrintStackPkgAlias
		imports[OtelGetStackImportPath] = OtelGetStackAliasPkg
	}
	for k, v := range imports {
		c += fmt.Sprintf("import %s \"%s\"\n", v, k)
	}

	// Assignments
	c += "func init() {\n"
	for _, assign := range assigns {
		c += assign
	}
	c += "}\n"

	f, err := util.WriteStringToFile(target, c)
	if err != nil {
		return err
	}
	dp.addGeneratedDep(f)
	shared.SaveDebugFile("", target)
	return err
}
func (dp *DepProcessor) setupOtelSDK(pkgName, target string) error {
	f, err := resource.CopyOtelSetupTo(pkgName, target)
	if err != nil {
		return fmt.Errorf("failed to copy otel setup sdk: %w", err)
	}
	dp.addGeneratedDep(f)
	shared.SaveDebugFile("", target)
	return err
}

// assembleInitCandidate assembles the candidate files that we may add init
// function to. The candidate files are the ones that have main or init
// function defined.
func assembleImportCandidates() ([]string, error) {
	candidates := make([]string, 0)
	found := false

	// Find from build arguments e.g. go build test_gorm_crud.go or go build cmd/app
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
	return candidates, nil
}

func (dp *DepProcessor) addExplicitImport(importPaths ...string) (err error) {
	// Find out where we should forcely import our init func
	if dp.importCandidates == nil {
		files, err := assembleImportCandidates()
		if err != nil {
			return fmt.Errorf("failed to assemble import candidates: %w", err)
		}
		dp.importCandidates = files
		if shared.Verbose {
			log.Printf("RuleImport candidates: %v", files)
		}
	}

	addImport := false
	for _, file := range dp.importCandidates {
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

	modFile, err := modfile.Parse("go.mod", []byte(data), nil)
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

func (dp *DepProcessor) setupRules() (err error) {
	err = os.MkdirAll(OtelRules, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %v: %w", OtelRules, err)
	}
	// Put otel_rule_*.go files into otel_rules
	err = dp.copyRules(OtelRules)
	if err != nil {
		return fmt.Errorf("failed to setup rules: %w", err)
	}
	// Put otel_setup_inst.go into otel_rules
	err = dp.initializeRules(OtelRules, filepath.Join(OtelRules, OtelSetupInst))
	if err != nil {
		return fmt.Errorf("failed to setup initiator: %w", err)
	}
	// Put otel_setup_sdk.go into otel_rules
	err = dp.setupOtelSDK(OtelRules, filepath.Join(OtelRules, OtelSetupSDK))
	if err != nil {
		return fmt.Errorf("failed to setup otel sdk: %w", err)
	}
	// Add implicit otel_rules import to introduce initialization side effect
	ruleImportPath, err := dp.getImportPathOf(OtelRules)
	if err != nil {
		return fmt.Errorf("failed to get import path: %w", err)
	}
	err = dp.addExplicitImport(ruleImportPath)
	if err != nil {
		return fmt.Errorf("failed to add rule import: %w", err)
	}
	return nil
}

func (dp *DepProcessor) addOtelImports() error {
	// We want to instrument otel-sdk itself, we done this by adding otel import
	// to the project, in this way, pkg/rules/otdk rules will always take effect.
	err := dp.addExplicitImport(
		OtelImportPath,
		OtelBaggageImportPath,
		OtelSdkTraceImportPath,
	)
	if err != nil {
		return fmt.Errorf("failed to add otel import: %w", err)
	}
	return nil
}

func (dp *DepProcessor) setupDeps() error {
	err := dp.addOtelImports()
	if err != nil {
		return fmt.Errorf("failed to add otel imports: %w", err)
	}

	compileCmds, err := getCompileCommands()
	if err != nil {
		return fmt.Errorf("failed to get compile commands: %w", err)
	}

	err = dp.matchRules(compileCmds)
	if err != nil {
		return fmt.Errorf("failed to find dependencies: %w", err)
	}

	err = dp.setupRules()
	if err != nil {
		return fmt.Errorf("failed to setup dependencies: %w", err)
	}
	return nil
}
