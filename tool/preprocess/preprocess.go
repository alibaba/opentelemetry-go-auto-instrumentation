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
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
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
	VendorDir          = "vendor"
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
		"v1.33.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace",
		"v1.33.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
		"v1.33.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
		"v1.33.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
		"v1.33.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp",
		"v1.33.0", true, false},
	{"go.opentelemetry.io/otel/exporters/prometheus",
		"v0.42.0", true, false},
	// otel contrib
	{"go.opentelemetry.io/contrib/instrumentation/runtime",
		"v0.58.0", false, false},
}

type DepProcessor struct {
	bundles         []*resource.RuleBundle // All dependent rule bundles
	backups         map[string]string
	localImportPath string
	sources         []string
	moduleName      string // Module name from go.mod
	modulePath      string // Where go.mod is located
	rule2Dir        map[*resource.InstFuncRule]string
	ruleCache       embed.FS
	goBuildCmd      []string
	vendorBuild     bool
}

func newDepProcessor() *DepProcessor {
	dp := &DepProcessor{
		bundles:         []*resource.RuleBundle{},
		backups:         map[string]string{},
		localImportPath: "",
		sources:         nil,
		rule2Dir:        map[*resource.InstFuncRule]string{},
		ruleCache:       pkg.ExportRuleCache(),
		vendorBuild:     false,
	}
	return dp
}

func (dp *DepProcessor) getGoModPath() string {
	util.Assert(dp.modulePath != "", "modulePath is empty")
	util.Assert(filepath.IsAbs(dp.modulePath), "modulePath is not absolute")
	return dp.modulePath
}

func (dp *DepProcessor) getGoModDir() string {
	return filepath.Dir(dp.getGoModPath())
}

func (dp *DepProcessor) getGoModName() string {
	return dp.moduleName
}

func (dp *DepProcessor) generatedOf(dir string) string {
	return filepath.Join(dp.getGoModDir(), dir)
}

// runCmdOutput runs the command and returns its standard output. dir specifies
// the working directory of the command. If dir is the empty string, run runs
// the command in the calling process's current directory.
func runCmdOutput(dir string, args ...string) (string, error) {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errc.New(errc.ErrRunCmd, string(out)).
			With("command", fmt.Sprintf("%v", args))
	}
	return string(out), nil
}

// Run runs the command and returns the combined standard output and standard
// error. dir specifies the working directory of the command. If dir is the
// empty string, run runs the command in the calling process's current directory.
func runCmdCombinedOutput(dir string, args ...string) (string, error) {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errc.New(errc.ErrRunCmd, string(out)).
			With("command", fmt.Sprintf("%v", args))
	}
	return string(out), nil
}

// Find go.mod from dir and its parent recursively
func findGoMod(dir string) (string, error) {
	for dir != "" {
		mod := filepath.Join(dir, util.GoModFile)
		if util.PathExists(mod) {
			return mod, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", errc.New(errc.ErrPreprocess, "cannot find go.mod")
}

func parseGoMod(gomod string) (*modfile.File, error) {
	data, err := util.ReadFile(gomod)
	if err != nil {
		return nil, err
	}
	modFile, err := modfile.Parse(util.GoModFile, []byte(data), nil)
	if err != nil {
		return nil, errc.New(errc.ErrParseCode, err.Error())
	}
	return modFile, nil
}

func (dp *DepProcessor) init() error {
	// There is a tricky, all arguments after the tool itself are saved for
	// later use, which means the subcommand "go build" are also  included
	dp.goBuildCmd = make([]string, len(os.Args)-1)
	copy(dp.goBuildCmd, os.Args[1:])
	util.AssertGoBuild(dp.goBuildCmd)

	// Find compiling module and package information from the build command
	pkgs, err := findModule(dp.goBuildCmd)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		util.Log("Find Go package %v", util.Jsonify(pkg))
		if pkg.GoFiles == nil {
			continue
		}
		dp.sources = append(dp.sources, pkg.GoFiles...)
		if pkg.Module != nil {
			// Best case, we find the module information from the package field
			util.Log("Find Go module %v", util.Jsonify(pkg.Module))
			util.Assert(pkg.Module.Path != "", "pkg.Module.Path is empty")
			util.Assert(pkg.Module.GoMod != "", "pkg.Module.GoMod is empty")
			dp.moduleName = pkg.Module.Path
			dp.modulePath = pkg.Module.GoMod
		} else {
			// If we cannot find the module information from the package field,
			// we try to find it from the go.mod file, where go.mod file is in
			// the same directory as the source file.
			util.Assert(pkg.Name != "", "pkg.Name is empty")
			if pkg.Name == "main" {
				gofile := pkg.GoFiles[0]
				gomod, err := findGoMod(filepath.Dir(gofile))
				if err != nil {
					return err
				}
				util.Assert(gomod != "", "gomod is empty")
				util.Assert(util.PathExists(gomod), "gomod does not exist")
				dp.modulePath = gomod
				// Get module name from go.mod file
				modfile, err := parseGoMod(gomod)
				if err != nil {
					return err
				}
				dp.moduleName = modfile.Module.Mod.Path
			}
		}
	}
	if len(dp.sources) == 0 {
		return errc.New(errc.ErrPreprocess, "no Go source files found")
	}
	if dp.moduleName == "" || dp.modulePath == "" {
		return errc.New(errc.ErrPreprocess, "cannot find compiled module")
	}

	util.Log("Found module %v in %v", dp.moduleName, dp.modulePath)
	util.Log("Found sources %v", dp.sources)

	// Check if the build mode
	ignoreVendor := false
	for _, arg := range dp.goBuildCmd {
		// -mod=mod and -mod=readonly tells the go command to ignore the vendor
		// directory. We should not use the vendor directory in this case.
		if strings.HasPrefix(arg, "-mod=mod") ||
			strings.HasPrefix(arg, "-mod=readonly") {
			dp.vendorBuild = false
			ignoreVendor = true
			break
		}
	}
	if !ignoreVendor {
		// FIXME: vendor directory name can be anything, but we assume it's "vendor"
		// for now
		vendor := filepath.Join(dp.getGoModDir(), VendorDir)
		dp.vendorBuild = util.PathExists(vendor)
	}
	util.Log("Vendor build: %v", dp.vendorBuild)

	// Register signal handler to catch up SIGINT/SIGTERM interrupt signals and
	// do necessary cleanup
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigc
		switch s {
		case syscall.SIGTERM, syscall.SIGINT:
			util.Log("Interrupted instrumentation, cleaning up")
			dp.postProcess()
		default:
		}
	}()

	return nil
}

func (dp *DepProcessor) postProcess() {
	util.GuaranteeInPreprocess()

	// Using -debug? Leave all changes for debugging
	if config.GetConf().Debug {
		return
	}

	// rm -rf otel_rules
	_ = os.RemoveAll(dp.generatedOf(OtelRules))

	// rm -rf otel_pkgdep
	_ = os.RemoveAll(dp.generatedOf(OtelPkgDep))

	// Restore everything we have modified during instrumentation
	_ = dp.restoreBackupFiles()
}

func (dp *DepProcessor) backupFile(origin string) error {
	util.GuaranteeInPreprocess()
	backup := filepath.Base(origin) + OtelBackupSuffix
	backup = util.GetLogPath(filepath.Join(OtelBackups, backup))
	err := os.MkdirAll(filepath.Dir(backup), 0777)
	if err != nil {
		return errc.New(errc.ErrMkdirAll, err.Error())
	}
	if _, exist := dp.backups[origin]; !exist {
		err = util.CopyFile(origin, backup)
		if err != nil {
			return err
		}
		dp.backups[origin] = backup
		util.Log("Backup %v", origin)
	} else if config.GetConf().Verbose {
		util.Log("Backup %v already exists", origin)
	}
	return nil
}

func (dp *DepProcessor) restoreBackupFiles() error {
	util.GuaranteeInPreprocess()
	for origin, backup := range dp.backups {
		err := util.CopyFile(backup, origin)
		if err != nil {
			return err
		}
		util.Log("Restore %v", origin)
	}
	return nil
}

func getCompileCommands() ([]string, error) {
	dryRunLog, err := os.Open(util.GetLogPath(DryRunLog))
	if err != nil {
		return nil, errc.New(errc.ErrOpenFile, err.Error())
	}
	defer func(dryRunLog *os.File) {
		err := dryRunLog.Close()
		if err != nil {
			util.Log("Failed to close dry run log file: %v", err)
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
		if util.IsCompileCommand(line) {
			line = strings.Trim(line, " ")
			compileCmds = append(compileCmds, line)
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, errc.New(errc.ErrParseCode, "cannot parse dry run log")
	}
	return compileCmds, nil
}

// $ go help packages
// Many commands apply to a set of packages:
//
//	go <action> [packages]
//
// Usually, [packages] is a list of import paths.
//
// An import path that is a rooted path or that begins with
// a . or .. element is interpreted as a file system path and
// denotes the package in that directory.
//
// Otherwise, the import path P denotes the package found in
// the directory DIR/src/P for some DIR listed in the GOPATH
// environment variable (For more details see: 'go help gopath').
//
// If no import paths are given, the action applies to the
// package in the current directory.
//
// There are four reserved names for paths that should not be used
// for packages to be built with the go tool:
//
// - "main" denotes the top-level package in a stand-alone executable.
//
// - "all" expands to all packages found in all the GOPATH
// trees. For example, 'go list all' lists all the packages on the local
// system. When using modules, "all" expands to all packages in
// the main module and their dependencies, including dependencies
// needed by tests of any of those.
//
// - "std" is like all but expands to just the packages in the standard
// Go library.
//
// - "cmd" expands to the Go repository's commands and their
// internal libraries.
//
// Import paths beginning with "cmd/" only match source code in
// the Go repository.
//
// An import path is a pattern if it includes one or more "..." wildcards,
// each of which can match any string, including the empty string and
// strings containing slashes. Such a pattern expands to all package
// directories found in the GOPATH trees with names matching the
// patterns.
//
// To make common patterns more convenient, there are two special cases.
// First, /... at the end of the pattern can match an empty string,
// so that net/... matches both net and packages in its subdirectories, like net/http.
// Second, any slash-separated pattern element containing a wildcard never
// participates in a match of the "vendor" element in the path of a vendored
// package, so that ./... does not match packages in subdirectories of
// ./vendor or ./mycode/vendor, but ./vendor/... and ./mycode/vendor/... do.
// Note, however, that a directory named vendor that itself contains code
// is not a vendored package: cmd/vendor would be a command named vendor,
// and the pattern cmd/... matches it.
// See golang.org/s/go15vendor for more about vendoring.
//
// An import path can also name a package to be downloaded from
// a remote repository. Run 'go help importpath' for details.
//
// Every package in a program must have a unique import path.
// By convention, this is arranged by starting each path with a
// unique prefix that belongs to you. For example, paths used
// internally at Google all begin with 'google', and paths
// denoting remote repositories begin with the path to the code,
// such as 'github.com/user/repo'.
//
// Packages in a program need not have unique package names,
// but there are two reserved package names with special meaning.
// The name main indicates a command, not a library.
// Commands are built into binaries and cannot be imported.
// The name documentation indicates documentation for
// a non-Go program in the directory. Files in package documentation
// are ignored by the go command.
//
// As a special case, if the package list is a list of .go files from a
// single directory, the command is applied to a single synthesized
// package made up of exactly those files, ignoring any build constraints
// in those files and ignoring any other files in the directory.
//
// Directory and file names that begin with "." or "_" are ignored
// by the go tool, as are directories named "testdata".

func tryLoadPackage(path string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		// Change it unless you know what you are doing
		Mode: packages.NeedModule | packages.NeedFiles | packages.NeedName,
	}

	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, errc.New(errc.ErrPreprocess, err.Error())
	}
	return pkgs, nil
}

func findModule(buildCmd []string) ([]*packages.Package, error) {
	candidates := make([]*packages.Package, 0)
	found := false

	// Find from build arguments e.g. go build test.go or go build cmd/app
	for i := len(buildCmd) - 1; i >= 0; i-- {
		buildArg := buildCmd[i]

		// Stop canary when we see a build flag or a "build" command
		if strings.HasPrefix("-", buildArg) || buildArg == "build" {
			break
		}

		// Trying to load package from the build argument, error is tolerated
		// because we dont know what the build argument is. One exception is
		// when we already found packages, in this case, we expect subsequent
		// build arguments are packages, so we should not tolerate any error.
		pkgs, err := tryLoadPackage(buildArg)
		if err != nil {
			if found {
				// If packages are already found, we expect subsequent build
				// arguments are packages, so we should not tolerate any error
				break
			}
			util.Log("Cannot load package from %v", buildArg)
			continue
		}
		for _, pkg := range pkgs {
			if pkg.Errors != nil {
				continue
			}
			found = true
			candidates = append(candidates, pkg)
		}
	}

	// If no import paths are given, the action applies to the package in the
	// current directory.
	if !found {
		pkgs, err := tryLoadPackage(".")
		if err != nil {
			return nil, err
		}
		for _, pkg := range pkgs {
			if pkg.Errors != nil {
				continue
			}
			candidates = append(candidates, pkg)
		}
	}
	if len(candidates) == 0 {
		return nil, errc.New(errc.ErrPreprocess, "no package found")
	}

	return candidates, nil
}

// Find out where we should forcely import our init func
func (dp *DepProcessor) addExplicitImport(importPaths ...string) (err error) {
	addImport := false
	for _, file := range dp.sources {
		if !util.IsGoFile(file) {
			continue
		}
		p := util.NewAstParser()
		astRoot, err := p.ParseFile(file, parser.ParseComments)
		if err != nil {
			return err
		}

		foundInit := util.FindFuncDecl(astRoot, FuncInit) != nil
		if !foundInit {
			foundMain := util.FindFuncDecl(astRoot, FuncMain) != nil
			if !foundMain {
				continue
			}
		}

		// Mark original line number of first import declaration
		for _, decl := range astRoot.Decls {
			decl, ok := decl.(*dst.GenDecl)
			if !ok {
				continue
			}
			if decl.Tok != token.IMPORT {
				continue
			}
			// The node is already tagged with line directive, dont tag it again
			generated := false
			if len(decl.Decs.Start) > 0 {
				for _, dec := range decl.Decs.Start {
					if strings.Contains(dec, "<generated>:1") {
						generated = true
						break
					}
				}
				if generated {
					break
				}
			}
			pos := p.FindPosition(decl)
			tag := fmt.Sprintf("//line %v", pos.String())
			decl.Decs.Before = dst.NewLine
			decl.Decs.Start.Append(tag)
			break
		}
		//  Prepend the generated import with a line directive to indicate
		//  that the generated code is not part of the original source code.
		util.Log("Add %v import to %v", importPaths, file)
		importDecl := util.AddImportForcely(astRoot, importPaths...)
		importDecl.Decs.Before = dst.NewLine
		importDecl.Decs.Start.Append("//line <generated>:1")
		addImport = true

		err = dp.backupFile(file)
		if err != nil {
			return err
		}
		_, err = util.WriteAstToFile(astRoot, filepath.Join(file))
		if err != nil {
			return err
		}
	}
	if !addImport {
		return errc.New(errc.ErrSetupRule, "no init or main function found")
	}
	return nil
}

func (dp *DepProcessor) findLocalImportPath() error {
	// Get absolute path of current working directory
	workingDir := dp.getGoModDir()
	// Get absolute path of go.mod directory
	projectDir := dp.getGoModDir()
	// Replace go.mod directory with module name
	moduleName := dp.getGoModName()
	// Replace all backslashes with slashes. The import path is different from
	// the file path, which should always use slashes.
	workingDir = filepath.ToSlash(workingDir)
	projectDir = filepath.ToSlash(projectDir)
	dp.localImportPath = strings.Replace(workingDir, projectDir, moduleName, 1)
	if config.GetConf().Verbose {
		util.Log("Find local import path: %v", dp.localImportPath)
	}
	return nil
}

func (dp *DepProcessor) getImportPathOf(dirName string) (string, error) {
	util.Assert(dirName != "", "dirName is empty")
	if dp.localImportPath == "" {
		err := dp.findLocalImportPath()
		if err != nil {
			return "", err
		}
	}
	// This is the import path in Go source code, not the file path, so we
	// should always use slashes.
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
		return err
	}
	return nil
}

// Note in this function, error is tolerated as we are best-effort to clean up
// any obsolete materials, but it's not fatal if we fail to do so.
func (dp *DepProcessor) preclean() {
	ruleImport, _ := dp.getImportPathOf(OtelRules)
	for _, file := range dp.sources {
		if !util.IsGoFile(file) {
			continue
		}
		astRoot, _ := util.ParseAstFromFile(file)
		if astRoot == nil {
			continue
		}
		if util.RemoveImport(astRoot, ruleImport) != nil {
			if config.GetConf().Verbose {
				util.Log("Remove obsolete import %v from %v",
					ruleImport, file)
			}
		}
		_, err := util.WriteAstToFile(astRoot, file)
		if err != nil {
			util.Log("Failed to write ast to %v: %v", file, err)
		}
	}
	// Clean otel_rules/otel_pkgdep directory
	if util.PathExists(dp.generatedOf(OtelRules)) {
		_ = os.RemoveAll(dp.generatedOf(OtelRules))
	}
	if util.PathExists(dp.generatedOf(OtelPkgDep)) {
		_ = os.RemoveAll(dp.generatedOf(OtelPkgDep))
	}
}

func (dp *DepProcessor) storeRuleBundles() error {
	err := resource.StoreRuleBundles(dp.bundles)
	if err != nil {
		return err
	}
	// No longer valid from now on
	dp.bundles = nil
	return nil
}

// runDryBuild runs a dry build to get all dependencies needed for the project.
func runDryBuild(goBuildCmd []string) error {
	dryRunLog, err := os.Create(util.GetLogPath(DryRunLog))
	if err != nil {
		return errc.New(errc.ErrCreateFile, err.Error())
	}
	// The full build command is: "go build -a -x -n  {...}"
	args := []string{"go", "build", "-a", "-x", "-n"}
	args = append(args, goBuildCmd[2:]...)
	util.AssertGoBuild(goBuildCmd)
	util.AssertGoBuild(args)

	// Run the dry build
	util.Log("Run dry build %v", args)
	cmd := exec.Command(args[0], args[1:]...)
	// This is a little anti-intuitive as the error message is not printed to
	// the stderr, instead it is printed to the stdout, only the build tool
	// knows the reason why.
	cmd.Stdout = os.Stdout
	cmd.Stderr = dryRunLog
	// @@Note that dir should not be set, as the dry build should be run in the
	// same directory as the original build command
	cmd.Dir = ""
	err = cmd.Run()
	if err != nil {
		return errc.New(errc.ErrRunCmd, err.Error()).
			With("command", fmt.Sprintf("%v", args))
	}
	return nil
}

func (dp *DepProcessor) runModTidy() error {
	out, err := runCmdCombinedOutput(dp.getGoModDir(),
		"go", "mod", "tidy")
	util.Log("Run go mod tidy: %v", out)
	return err
}

func (dp *DepProcessor) runModVendor() error {
	out, err := runCmdCombinedOutput(dp.getGoModDir(),
		"go", "mod", "vendor")
	util.Log("Run go mod vendor: %v", out)
	return err
}

func (dp *DepProcessor) runGoGet(dep string) error {
	out, err := runCmdCombinedOutput(dp.getGoModDir(),
		"go", "get", dep)
	util.Log("Run go get %v: %v", dep, out)
	return err
}

func (dp *DepProcessor) runGoModDownload(path string) error {
	out, err := runCmdCombinedOutput(dp.getGoModDir(),
		"go", "mod", "download", path)
	util.Log("Run go mod download %v: %v", path, out)
	return err
}

func (dp *DepProcessor) runGoModEdit(require string) error {
	out, err := runCmdCombinedOutput(dp.getGoModDir(),
		"go", "mod", "edit", "-require="+require)
	util.Log("Run go mod edit %v: %v", require, out)
	return err
}

func nullDevice() string {
	if runtime.GOOS == "windows" {
		return "NUL"
	}
	return "/dev/null"
}

func runBuildWithToolexec(goBuildCmd []string) error {
	exe, err := os.Executable()
	if err != nil {
		return errc.New(errc.ErrGetExecutable, err.Error())
	}
	args := []string{
		"go",
		"build",
		// Add remix subcommand to tell the tool this is toolexec mode
		"-toolexec=" + exe + " " + CompileRemix,
	}

	// Leave the temporary compilation directory
	args = append(args, util.BuildWork)

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
		util.Log("Run go build with args %v in toolexec mode", args)
	}
	util.AssertGoBuild(args)
	// @@ Note that we should not set the working directory here, as the build
	// with toolexec should be run in the same directory as the original build
	// command
	out, err := runCmdCombinedOutput("", args...)
	util.Log("Run go build with toolexec: %v", out)
	return err
}

func (dp *DepProcessor) fetchDep(path string) error {
	err := dp.runGoModDownload(path)
	if err != nil {
		return err
	}
	err = dp.runGoModEdit(path)
	if err != nil {
		return err
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
			util.Log("Pin dependency version %v@%v", p, v)
		}
		err := dp.fetchDep(p + "@" + v)
		if err != nil {
			if dep.fallible {
				util.Log("Failed to pin dependency %v: %v", p, err)
				continue
			}
			return err
		}
	}
	return nil
}

func precheck() error {
	// Check if the project is modularized
	go11module := os.Getenv("GO111MODULE")
	if go11module == "off" {
		return errc.New(errc.ErrNotModularized, "GO111MODULE is off")
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
	gomodDir := dp.getGoModDir()
	files := []string{}
	files = append(files, filepath.Join(gomodDir, util.GoModFile))
	files = append(files, filepath.Join(gomodDir, util.GoSumFile))
	files = append(files, filepath.Join(gomodDir, util.GoWorkSumFile))
	for _, file := range files {
		if util.PathExists(file) {
			err := dp.backupFile(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (dp *DepProcessor) saveDebugFiles() {
	dir := filepath.Join(util.GetTempBuildDir(), OtelRules)
	err := os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		util.CopyDir(dp.generatedOf(OtelRules), dir)
	}
	dir = filepath.Join(util.GetTempBuildDir(), OtelUser)
	err = os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		for origin := range dp.backups {
			util.CopyFile(origin, filepath.Join(dir, filepath.Base(origin)))
		}
	}
}

func (dp *DepProcessor) setupDeps() error {
	// Pre-clean before processing in case of any obsolete materials left
	dp.preclean()

	err := dp.addOtelImports()
	if err != nil {
		return err
	}

	// Pinning otel version in go.mod
	err = dp.pinDepVersion()
	if err != nil {
		return err
	}

	// Run go mod tidy first to fetch all dependencies
	err = dp.runModTidy()
	if err != nil {
		return err
	}

	if dp.vendorBuild {
		err = dp.runModVendor()
		if err != nil {
			return err
		}
	}

	// Run dry build to the build blueprint
	err = runDryBuild(dp.goBuildCmd)
	if err != nil {
		// Tell us more about what happened in the dry run
		errLog, _ := util.ReadFile(util.GetLogPath(DryRunLog))
		err = errc.Adhere(err, "reason", errLog)
		return err
	}

	// Find compile commands from dry run log
	compileCmds, err := getCompileCommands()
	if err != nil {
		return err
	}

	err = dp.copyPkgDep()
	if err != nil {
		return err
	}

	// Find used rules according to compile commands
	err = dp.matchRules(compileCmds)
	if err != nil {
		return err
	}

	err = dp.fetchRules()
	if err != nil {
		return err
	}

	// Setup rules according to compile commands
	err = dp.setupRules()
	if err != nil {
		return err
	}

	err = dp.replaceOtelImports()
	if err != nil {
		return err
	}

	// Save matched rules into file, from this point on, we no longer modify
	// the rules
	err = dp.storeRuleBundles()
	if err != nil {
		return err
	}
	return nil
}

func Preprocess() error {
	// Make sure the project is modularized otherwise we cannot proceed
	err := precheck()
	if err != nil {
		return err
	}

	dp := newDepProcessor()

	err = dp.init()
	if err != nil {
		return err
	}

	defer func() { dp.postProcess() }()
	{
		defer util.PhaseTimer("Preprocess")()

		// Backup go.mod as we are likely modifing it later
		err = dp.backupMod()
		if err != nil {
			return err
		}

		// Run a dry build to get all dependencies needed for the project
		// Match the dependencies with available rules and prepare them
		// for the actual instrumentation
		err = dp.setupDeps()
		if err != nil {
			return err
		}

		// Pinning dependencies version in go.mod
		err = dp.pinDepVersion()
		if err != nil {
			return err
		}

		// Run go mod tidy to fetch dependencies
		err = dp.runModTidy()
		if err != nil {
			return err
		}

		if dp.vendorBuild {
			err = dp.runModVendor()
			if err != nil {
				return err
			}
		}

		// 	// Retain otel rules and modified user files for debugging
		dp.saveDebugFiles()
	}

	{
		defer util.PhaseTimer("Instrument")()

		// Run go build with toolexec to start instrumentation
		err = runBuildWithToolexec(dp.goBuildCmd)
		if err != nil {
			return err
		}
	}
	util.Log("Build completed successfully")
	return nil
}
