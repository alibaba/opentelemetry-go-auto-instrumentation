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
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

// -----------------------------------------------------------------------------
// Preprocess
//
// The preprocess package is used to preprocess the source code before the actual
// instrumentation. Instrumentation rules may introduces additional dependencies
// that are not present in original source code. The preprocess is responsible
// for preparing these dependencies in advance.

const (
	OtelImportHolder   = "otel_import_holder.go"
	OtelSetupInst      = "otel_setup_inst.go"
	OtelSetupSDK       = "otel_setup_sdk.go"
	OtelRules          = "otel_rules"
	OtelUser           = "otel_user"
	OtelRuleCache      = "rule_cache"
	OtelBackups        = "backups"
	OtelBackupSuffix   = ".bk"
	FuncMain           = "main"
	FuncInit           = "init"
	DryRunLog          = "dry_run.log"
	CompileRemix       = "remix"
	ReorderLocalPrefix = "<REORDER>"
	ReorderInitFile    = "reorder_init.go"
	StdRulesPrefix     = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/"
	StdRulesPath       = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules"
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

var fixedDeps = []struct {
	dep, version string
	addImport    bool
	fallible     bool
}{
	// otel sdk
	{"go.opentelemetry.io/otel",
		"v1.32.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace",
		"v1.32.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
		"v1.32.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
		"v1.32.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
		"v1.32.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp",
		"v1.32.0", true, false},
	// otel contrib
	{"go.opentelemetry.io/contrib/instrumentation/runtime",
		"v0.57.0", false, false},
	// otel itself
	// {"github.com/alibaba/opentelemetry-go-auto-instrumentation",
	// "v0.3.0", false, true},
}

type DepProcessor struct {
	bundles         []*resource.RuleBundle // All dependent rule bundles
	backups         map[string]string
	localImportPath string
	importHolder    string
	rule2Dir        map[*resource.InstFuncRule]string
	ruleCache       embed.FS
	goBuildCmd      []string
}

func newDepProcessor() *DepProcessor {
	dp := &DepProcessor{
		bundles:         []*resource.RuleBundle{},
		backups:         map[string]string{},
		localImportPath: "",
		importHolder:    "",
		rule2Dir:        map[*resource.InstFuncRule]string{},
		ruleCache:       pkg.ExportRuleCache(),
	}
	return dp
}

func (dp *DepProcessor) init() error {
	// Pre-clean before processing in case of any obsolete materials left
	dp.cleanup()

	// Parse the go build command
	// There is a tricky, all arguments after the tool itself are saved for
	// later use, which means the subcommand "go build" are also included
	dp.goBuildCmd = make([]string, len(os.Args)-1)
	copy(dp.goBuildCmd, os.Args[1:])
	shared.AssertGoBuild(dp.goBuildCmd)

	// Initialize the import holder file, all explicit imports will be added
	// to this file
	util.Assert(dp.importHolder == "", "already initialized")
	util.Assert(len(dp.goBuildCmd) >= 2, "invalid go build command")
	pkgs, err := loadPackages(dp.goBuildCmd[2:]...)
	if err != nil {
		return fmt.Errorf("failed to load packages: %w", err)
	}
	mainPkg, err := findMainPackage(pkgs)
	if err != nil {
		return fmt.Errorf("failed to find main package: %w", err)
	}
	if len(mainPkg.GoFiles) == 0 {
		return fmt.Errorf("no go files in main package")
	}
	dir := filepath.Dir(mainPkg.GoFiles[0])
	dp.importHolder = filepath.Join(dir, OtelImportHolder)
	_, err = util.WriteFile(dp.importHolder, fmt.Sprintf("package main\n"))
	if err != nil {
		return fmt.Errorf("failed to write import holder: %w", err)
	}
	log.Printf("Create %v\n", dp.importHolder)

	// Register signal handler to catch up SIGINT/SIGTERM interrupt signals and
	// do necessary cleanup
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigc
		switch s {
		case syscall.SIGTERM, syscall.SIGINT:
			log.Printf("Interrupted instrumentation, cleaning up")
			dp.postProcess()
		default:
		}
	}()
	return nil
}

func (dp *DepProcessor) cleanup() {
	// rm -rf import holder
	if dp.importHolder != "" {
		_ = os.Remove(dp.importHolder)
	}

	// rm -rf otel_rules
	_ = os.RemoveAll(OtelRules)

	// rm -rf otel_pkgdep
	_ = os.RemoveAll(OtelPkgDep)
}

// postProcess is called after the instrumentation is done. It cleans up all
// temporary files and restores all changes made during the preprocessing.
func (dp *DepProcessor) postProcess() {
	shared.GuaranteeInPreprocess()
	// Clean build cache as we may instrument some std packages(e.g. runtime)
	// TODO: fine-grained cache cleanup
	// err := util.RunCmd("go", "clean", "-cache")
	// if err != nil {
	// 	log.Fatalf("failed to clean cache: %v", err)
	// }

	// Using -debug? Leave all changes for debugging
	if config.GetConf().Debug {
		return
	}

	dp.cleanup()

	// Restore everything we have modified during instrumentation
	err := dp.restoreBackupFiles()
	if err != nil {
		log.Fatalf("failed to restore: %v", err)
	}
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
	} else if config.GetConf().Verbose {
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

func loadPackages(patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("failed to load packages")
	}
	return pkgs, nil
}

func findMainPackage(pkgs []*packages.Package) (*packages.Package, error) {
	for _, pkg := range pkgs {
		for _, gofile := range pkg.GoFiles {
			astRoot, err := shared.ParseAstFromFile(gofile)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ast from file: %w", err)
			}
			if shared.FindFuncDecl(astRoot, FuncMain) != nil {
				return pkg, nil
			}
		}
	}
	return nil, fmt.Errorf("main package not found")
}

func (dp *DepProcessor) addExplicitImport(importPaths ...string) (err error) {
	for _, importPath := range importPaths {
		_, err = util.AppendFile(dp.importHolder,
			fmt.Sprintf("import _ %q\n", importPath),
		)
		if err != nil {
			return fmt.Errorf("failed to append import: %w", err)
		}
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
	// Note that import path is different from file path, it should be presented
	// in slash-separated form on all platforms
	workingDir = filepath.ToSlash(workingDir)
	projectDir = filepath.ToSlash(projectDir)
	dp.localImportPath = strings.Replace(workingDir, projectDir, moduleName, 1)
	if config.GetConf().Verbose {
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

func (dp *DepProcessor) storeRuleBundles() error {
	err := resource.StoreRuleBundles(dp.bundles)
	if err != nil {
		return fmt.Errorf("failed to store rule bundles: %w", err)
	}
	// No longer valid from now on
	dp.bundles = nil
	return nil
}

// runDryBuild runs a dry build to get all dependencies needed for the project.
func runDryBuild(goBuildCmd []string) error {
	dryRunLog, err := os.Create(shared.GetLogPath(DryRunLog))
	if err != nil {
		return err
	}
	// The full build command is: "go build -a -x -n  {...}"
	args := []string{"go", "build", "-a", "-x", "-n"}
	args = append(args, goBuildCmd...)
	args = util.StringDedup(args)
	shared.AssertGoBuild(args)

	// Run the dry build
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = dryRunLog
	cmd.Stderr = dryRunLog
	return cmd.Run()
}

func runModTidy() error {
	return util.RunCmd("go", "mod", "tidy")
}

func runGoGet(dep string) error {
	return util.RunCmd("go", "get", dep)
}

func runGoModDownload(path string) error {
	return util.RunCmd("go", "mod", "download", path)
}

func runGoModEdit(require string) error {
	return util.RunCmd("go", "mod", "edit", "-require="+require)
}

func runCleanCache() error {
	return util.RunCmd("go", "clean", "-cache")
}

func nullDevice() string {
	if runtime.GOOS == "windows" {
		return "NUL"
	}
	return "/dev/null"
}

func runBuildWithToolexec(goBuildCmd []string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	args := []string{
		"go",
		"build",
		// Add remix subcommand to tell the tool this is toolexec mode
		"-toolexec=" + exe + " " + CompileRemix,
	}

	// Leave the temporary compilation directory
	args = append(args, "-work")

	// Force rebuilding
	args = append(args, "-a")

	if config.GetConf().Debug {
		// Disable compiler optimizations for debugging mode
		args = append(args, "-gcflags=all=-N -l")
	}

	// Append additional build arguments provided by the user
	args = append(args, goBuildCmd[2:]...)

	if config.GetConf().Restore {
		// Dont generate any compiled binary when using -restore
		args = append(args, "-o")
		args = append(args, nullDevice())
	}

	if config.GetConf().Verbose {
		log.Printf("Run go build with args %v in toolexec mode", args)
	}
	args = util.StringDedup(args)
	shared.AssertGoBuild(args)
	return util.RunCmdOutput(args...)
}

func fetchDep(path string) error {
	err := runGoModDownload(path)
	if err != nil {
		return fmt.Errorf("failed to fetch dependency %v: %w", path, err)
	}
	err = runGoModEdit(path)
	if err != nil {
		return fmt.Errorf("failed to edit go.mod: %w", err)
	}
	return nil
}

// We want to fetch otel dependencies in a fixed version instead of the latest
// version, so we need to pin the version in go.mod. All used otel dependencies
// should be listed and pinned here, because go mod tidy will fetch the latest
// version even if we have pinned some of them.
// Users will import github.com/alibaba/opentelemetry-go-auto-instrumentation
// dependency while using otel to use the inst-api and inst-semconv package
// We also need to pin its version to let the users use the fixed version
func (dp *DepProcessor) pinDepVersion() error {
	// otel related sdk dependencies
	for _, dep := range fixedDeps {
		p := dep.dep
		v := dep.version
		if config.GetConf().Verbose {
			log.Printf("Pin dependency version %v@%v", p, v)
		}
		err := fetchDep(p + "@" + v)
		if err != nil {
			if dep.fallible {
				log.Printf("Failed to pin dependency %v: %v", p, err)
				continue
			}
			return fmt.Errorf("failed to pin dependency %v: %w", dep, err)
		}
	}
	return nil
}

func precheck() error {
	// Check if the project is modularized
	go11module := os.Getenv("GO111MODULE")
	if go11module == "off" {
		return fmt.Errorf("GO111MODULE is set to off")
	}
	found, err := shared.IsExistGoMod()
	if !found {
		return fmt.Errorf("go.mod not found %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to check go.mod: %w", err)
	}
	// Check if the project is build with vendor mode
	projRoot, err := shared.GetGoModDir()
	if err != nil {
		return fmt.Errorf("failed to get project root: %w", err)
	}
	vendor := filepath.Join(projRoot, shared.VendorDir)
	if exist, _ := util.PathExists(vendor); exist {
		return fmt.Errorf("vendor mode is not supported")
	}

	// Check if the build arguments is sane
	if len(os.Args) < 3 {
		config.PrintVersion()
		os.Exit(0)
	}
	if !strings.Contains(os.Args[1], "go") {
		config.PrintVersion()
		os.Exit(0)
	}
	if os.Args[2] != "build" {
		config.PrintVersion()
		os.Exit(0)
	}
	return nil
}

func (dp *DepProcessor) backupMod() error {
	gomodDir, err := shared.GetGoModDir()
	if err != nil {
		return fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	files := []string{}
	files = append(files, filepath.Join(gomodDir, shared.GoModFile))
	files = append(files, filepath.Join(gomodDir, shared.GoSumFile))
	files = append(files, filepath.Join(gomodDir, shared.GoWorkSumFile))
	for _, file := range files {
		if exist, _ := util.PathExists(file); exist {
			err = dp.backupFile(file)
			if err != nil {
				return fmt.Errorf("failed to backup %s: %w", file, err)
			}
		}
	}
	return nil
}

func (dp *DepProcessor) saveDebugFiles() {
	dir := filepath.Join(shared.GetTempBuildDir(), OtelRules)
	err := os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		util.CopyDir(OtelRules, dir)
	}
	dir = filepath.Join(shared.GetTempBuildDir(), OtelUser)
	err = os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		for origin := range dp.backups {
			util.CopyFile(origin, filepath.Join(dir, filepath.Base(origin)))
		}
	}
}

func (dp *DepProcessor) setupDeps() error {
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

	// Run dry build to the build blueprint
	err = runDryBuild(dp.goBuildCmd)
	if err != nil {
		// Tell us more about what happened in the dry run
		errLog, _ := util.ReadFile(shared.GetLogPath(DryRunLog))
		return fmt.Errorf("failed to run dry build: %w\n%v", err, errLog)
	}

	// Find compile commands from dry run log
	compileCmds, err := getCompileCommands()
	if err != nil {
		return fmt.Errorf("failed to get compile commands: %w", err)
	}

	err = dp.copyPkgDep()
	if err != nil {
		return fmt.Errorf("failed to copy pkgdep: %w", err)
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

	err = dp.replaceOtelImports()
	if err != nil {
		return fmt.Errorf("failed to replace otel imports: %w", err)
	}

	// Save matched rules into file, from this point on, we no longer modify
	// the rules
	err = dp.storeRuleBundles()
	if err != nil {
		return fmt.Errorf("failed to store rule bundles: %w", err)
	}
	return nil
}

func Preprocess() error {
	// Make sure the project is modularized otherwise we cannot proceed
	err := precheck()
	if err != nil {
		return fmt.Errorf("not modularized project: %w", err)
	}

	dp := newDepProcessor()
	err = dp.init()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}
	defer func() { dp.postProcess() }()
	{
		defer util.PhaseTimer("Preprocess")()

		// Backup go.mod as we are likely modifing it later
		err = dp.backupMod()
		if err != nil {
			return fmt.Errorf("failed to backup go.mod: %w", err)
		}

		// Run a dry build to get all dependencies needed for the project
		// Match the dependencies with available rules and prepare them
		// for the actual instrumentation
		err = dp.setupDeps()
		if err != nil {
			return fmt.Errorf("failed to setup prerequisites: %w", err)
		}

		// Pinning dependencies version in go.mod
		err = dp.pinDepVersion()
		if err != nil {
			return fmt.Errorf("failed to update otel: %w", err)
		}

		// Run go mod tidy to fetch dependencies
		err = runModTidy()
		if err != nil {
			return fmt.Errorf("failed to run mod tidy: %w", err)
		}

		// 	// Retain otel rules and modified user files for debugging
		dp.saveDebugFiles()
	}

	{
		defer util.PhaseTimer("Instrument")()

		// Run go build with toolexec to start instrumentation
		out, err := runBuildWithToolexec(dp.goBuildCmd)
		if err != nil {
			return fmt.Errorf("failed to run go toolexec build: %w\n%s",
				err, out)
		} else {
			log.Printf("CompileRemix: %s", out)
		}
	}
	return nil
}
