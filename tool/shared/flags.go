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

package shared

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	TempBuildDir    = ".otel-build"
	VendorDir       = "vendor"
	BuildModeVendor = "-mod=vendor"
	BuildModeMod    = "-mod=mod"
)

// InToolexec true means this tool is being invoked in the go build process.
// This flag should not be set manually by users.
var InToolexec bool

// DebugLog true means debug log is enabled.
var DebugLog = false

// Verbose true means print verbose log.
var Verbose = true

// DisableRules enable rules by name(* for all, comma separated names).
var DisableRules = "testrule"

// Debug true means debug mode.
var Debug = false

// Restore true means restore all instrumentations.
var Restore = false

// BuildArgs are the arguments to pass to the go build command.
var BuildArgs []string

// Version
var PrintVersion = false

// The following flags should be shared across preprocess and instrument.
const (
	WorkingDirEnv   = "OTEL_WORKING_DIRECTORY"
	DebugLogEnv     = "OTEL_DEBUG_TO_FILE"
	DisableRulesEnv = "OTEL_DISABLE_RULES"
	VerboseEnv      = "OTEL_VERBOSE"
)

// This is the version of the tool, which will be printed when the -version flag
// is passed. This value is specified by the build system.
var TheVersion = "1.0.0"

var TheName = "otelbuild"

func PrintTheVersion() {
	fmt.Printf("%s version %s\n", TheName, TheVersion)
}

func ParseOptions() {
	// Parse flags from command-line arguments
	flag.BoolVar(&InToolexec, "in-toolexec", false, "Run in toolexec mode")
	flag.BoolVar(&DebugLog, "debuglog", false, "Print debug log to file")
	flag.BoolVar(&Verbose, "verbose", false, "Print verbose log")
	flag.StringVar(&DisableRules, "disablerules", "testrule", "Enable rules by name, * for all, comma separated names")
	flag.BoolVar(&Debug, "debug", false, "Enable debug mode, leave temporary files for debugging")
	flag.BoolVar(&Restore, "restore", false, "Restore all instrumentations")
	flag.BoolVar(&PrintVersion, "version", false, "Print version")
	flag.Parse()

	// Any non-flag command-line arguments behind "--" separator will be treated
	// as build arguments and transparently passed to the go build command.
	BuildArgs = flag.Args()
}

func SetBuildMode() error {
	// We can't brutely always add -mod=mod here because -mod may only be set to
	// readonly or vendor when in workspace mode. We need to check if provided
	// -mod is vendor or vendor directory exists, then we set -mod=mod mode.
	// For all other cases, we just leave it as is.

	// Check if -mod=vendor is set, replace it with -mod=mod
	const buildModeVendor = "-mod=vendor"
	const buildModePrefix = "-mod"
	for i, arg := range BuildArgs {
		// -mod=vendor?
		if arg == buildModeVendor {
			BuildArgs[i] = BuildModeMod
			return nil
		}
		// -mod vendor?
		if arg == buildModePrefix {
			if len(BuildArgs) > i+1 {
				arg1 := BuildArgs[i+1]
				if arg1 == "vendor" {
					BuildArgs[i+1] = "mod"
					return nil
				}
			}
		}
	}

	// Check if vendor directory exists, explicitly set -mod=mod
	gomodDir, err := GetGoModDir()
	if err != nil {
		return fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	vendor := filepath.Join(gomodDir, VendorDir)
	exist, err := util.PathExists(vendor)
	if err != nil {
		return fmt.Errorf("failed to check vendor directory: %w", err)
	}
	if exist {
		BuildArgs = append([]string{BuildModeMod}, BuildArgs...)
	}
	return nil
}

func InitOptions() (err error) {
	if InToolexec {
		// Inherit options from environment variables
		wd := os.Getenv(WorkingDirEnv)
		if len(wd) == 0 {
			return fmt.Errorf("cannot find working directory")
		}

		DebugLog, _ = strconv.ParseBool(os.Getenv(DebugLogEnv))
		Verbose, _ = strconv.ParseBool(os.Getenv(VerboseEnv))
		DisableRules = os.Getenv(DisableRulesEnv)
	} else {
		wd, err := filepath.Abs(TempBuildDir)
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		// Otherwise, set environment variables for further go toolexec build
		if err = os.Setenv(WorkingDirEnv, wd); err != nil {
			return fmt.Errorf("failed to set working directory: %w", err)
		}
		if err = os.Setenv(DebugLogEnv, strconv.FormatBool(DebugLog)); err != nil {
			return fmt.Errorf("failed to set debug log flag: %w", err)
		}
		if err = os.Setenv(DisableRulesEnv, DisableRules); err != nil {
			return fmt.Errorf("failed to set use rules flag: %w", err)
		}
		if err = os.Setenv(VerboseEnv, strconv.FormatBool(Verbose)); err != nil {
			return fmt.Errorf("failed to set use rules flag: %w", err)
		}
	}
	return SetBuildMode()
}
