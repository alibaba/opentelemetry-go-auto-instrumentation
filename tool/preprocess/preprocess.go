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
	"runtime"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
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

const FixedOtelDepVersion = "v1.28.0"

var fixedOtelDeps = []string{
	"go.opentelemetry.io/otel",
	"go.opentelemetry.io/otel/sdk",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
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

	// Avoid conflicts between args and flags with -- symbol
	// According to go1.18.10/src/flag/flag.go#parseOne
	if len(shared.BuildArgs) > 0 {
		if shared.Restore {
			return "", fmt.Errorf("args are not allowed when -restore presents")
		}
		// Append additional build arguments provided by the user
		args = append(args, shared.BuildArgs...)
	}

	if shared.Restore {
		// Dont generate any compiled binary when using -restore
		args = append(args, "-o")
		args = append(args, nullDevice())
	}

	if shared.Verbose {
		log.Printf("Run go build with args %v in toolexec mode", args)
	}
	return util.RunCmdWithOutput(append([]string{"go"}, args...)...)
}

func (dp *DepProcessor) pinDepVersion() error {
	// This should be done before running go mod tidy, because we may relies on
	// some packages that only presents in our specified version, running go mod
	// tidy will report error since it nevertheless pulls the latest version,
	// which is not what we want.
	for _, ruleHash := range dp.funcRules {
		rule := resource.FindFuncRuleByHash(ruleHash)
		for _, dep := range rule.PackageDeps {
			log.Printf("Pin dependency version %v ", dep)
			err := runGoGet(dep)
			if err != nil {
				return fmt.Errorf("failed to pin dependency %v: %w", dep, err)
			}
		}
	}
	return nil
}

// We want to fetch otel dependencies in a fixed version instead of the latest
// version, so we need to pin the version in go.mod. All used otel dependencies
// should be listed and pinned here, because go mod tidy  will fetch the latest
// version even if we have pinned some of them.
func (dp *DepProcessor) pinOtelVersion() error {
	for _, dep := range fixedOtelDeps {
		log.Printf("Pin otel dependency version %v ", dep)
		err := runGoGet(dep + "@" + FixedOtelDepVersion)
		if err != nil {
			return fmt.Errorf("failed to pin otel dependency %v: %w", dep, err)
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

func (dp *DepProcessor) backupGoMod() error {
	gomod, err := shared.GetGoModPath()
	if err != nil {
		return fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	if gomod == "" {
		return fmt.Errorf("failed to get go.mod directory")
	}
	err = dp.backupFile(gomod)
	if err != nil {
		return fmt.Errorf("failed to backup go.mod: %w", err)
	}
	return nil
}

func Preprocess() error {
	err := checkModularized()
	if err != nil {
		return fmt.Errorf("not modularized project: %w", err)
	}

	dp := newDepProcessor()
	// Aware of both normal exit and interrupt signal, clean up temporary files
	defer func() { dp.postProcess() }()
	dp.catchSignal()
	{
		start := time.Now()
		// Backup go.mod as we are likely modifing it later
		err = dp.backupGoMod()
		if err != nil {
			return fmt.Errorf("failed to backup go.mod: %w", err)
		}

		// Find rule dependencies according to compile commands
		err = dp.setupDeps()
		if err != nil {
			return fmt.Errorf("failed to setup prerequisites: %w", err)
		}
		log.Printf("Setup rules took %v", time.Since(start))
		start = time.Now()

		// Pinning dependencies version in go.mod
		err = dp.pinDepVersion()
		if err != nil {
			return fmt.Errorf("failed to update dependencies: %w", err)
		}

		// Pinning otel version in go.mod
		err = dp.pinOtelVersion()
		if err != nil {
			return fmt.Errorf("failed to update otel: %w", err)
		}

		// Run go mod tidy to fetch dependencies
		err = runModTidy()
		if err != nil {
			return fmt.Errorf("failed to run mod tidy: %w", err)
		}

		log.Printf("Preprocess took %v", time.Since(start))
	}

	{
		start := time.Now()
		// Run go build with toolexec to start instrumentation
		out, err := runBuildWithToolexec()
		if err != nil {
			return fmt.Errorf("failed to run go toolexec build: %w\n%s",
				err, out)
		}
		log.Printf("Instrument took %v", time.Since(start))
	}
	return nil
}
