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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

// -----------------------------------------------------------------------------
// Preprocess
//
// The preprocess package is used to preprocess the source code before the actual
// instrumentation. Instrumentation rules may introduces additional dependencies
// that are not present in original source code. The preprocess is responsible
// for preparing these dependencies in advance.

const (
	DryRunLog = "dry_run.log"
)

var fixedDeps = []struct {
	dep, version string
	addImport    bool
	fallible     bool
}{
	// otel sdk
	{"go.opentelemetry.io/otel",
		"v1.30.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace",
		"v1.30.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
		"v1.30.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
		"v1.30.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
		"v1.30.0", true, false},
	{"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp",
		"v1.30.0", true, false},
	// otel contrib
	{"go.opentelemetry.io/contrib/instrumentation/runtime",
		"v0.55.0", false, false},
	// otelbuild itself
	{"github.com/alibaba/opentelemetry-go-auto-instrumentation",
		"v0.2.0.dev", false, true},
}

// runDryBuild runs a dry build to get all dependencies needed for the project.
func runDryBuild() error {
	dryRunLog, err := os.Create(shared.GetLogPath(DryRunLog))
	if err != nil {
		return err
	}
	// The full build command is: "go build -a -x -n {BuildArgs...}"
	args := append([]string{"build", "-a", "-x", "-n"}, shared.BuildArgs...)
	cmd := exec.Command("go", args...)
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

func runCleanCache() error {
	return util.RunCmd("go", "clean", "-cache")
}

func nullDevice() string {
	if runtime.GOOS == "windows" {
		return "NUL"
	}
	return "/dev/null"
}

func runBuildWithToolexec() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	args := []string{
		"build",
		"-toolexec=" + exe + " -in-toolexec",
	}

	// Leave the temporary compilation directory
	args = append(args, "-work")

	// Force rebuilding
	args = append(args, "-a")

	if shared.Debug {
		// Disable compiler optimizations for debugging mode
		args = append(args, "-gcflags=all=-N -l")
	}

	// Append additional build arguments provided by the user
	args = append(args, shared.BuildArgs...)

	if shared.Restore {
		// Dont generate any compiled binary when using -restore
		args = append(args, "-o")
		args = append(args, nullDevice())
	}

	if shared.Verbose {
		log.Printf("Run go build with args %v in toolexec mode", args)
	}
	return util.RunCmdOutput(append([]string{"go"}, args...)...)
}

// We want to fetch otel dependencies in a fixed version instead of the latest
// version, so we need to pin the version in go.mod. All used otel dependencies
// should be listed and pinned here, because go mod tidy will fetch the latest
// version even if we have pinned some of them.
// Users will import github.com/alibaba/opentelemetry-go-auto-instrumentation
// dependency while using otelbuild to use the inst-api and inst-semconv package.
// We also need to pin its version to let the users use the fixed version
func (dp *DepProcessor) pinDepVersion() error {
	// otel related sdk dependencies
	for _, dep := range fixedDeps {
		p := dep.dep
		v := dep.version
		log.Printf("Pin dependency version %v@%v", p, v)
		err := runGoGet(p + "@" + v)
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

func checkModularized() error {
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

func Preprocess() error {
	// Make sure the project is modularized otherwise we cannot proceed
	err := checkModularized()
	if err != nil {
		return fmt.Errorf("not modularized project: %w", err)
	}

	dp := newDepProcessor()
	defer func() { dp.postProcess() }()
	dp.catchSignal()

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
		out, err := runBuildWithToolexec()
		if err != nil {
			return fmt.Errorf("failed to run go toolexec build: %w\n%s",
				err, out)
		} else {
			log.Printf("CompileRemix: %s", out)
		}
	}
	return nil
}
