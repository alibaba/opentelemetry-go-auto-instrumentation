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
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	TempBuildDir    = ".otel-build"
	VendorDir       = "vendor"
	BuildModeVendor = "-mod=vendor"
	BuildModeMod    = "-mod=mod"
)

// RuleFile is the file name of the rule file.
var RuleJsonFiles = ""

// InToolexec true means this tool is being invoked in the go build process.
// This flag should not be set manually by users.
var InToolexec bool

// DebugLog true means debug log is enabled.
var DebugLog = false

// Verbose true means print verbose log.
var Verbose = true

// Debug true means debug mode.
var Debug = false

// Restore true means restore all instrumentations.
var Restore = false

// GoBuildCmd are the original go build command.
var GoBuildCmd []string

// Version
var PrintVersion = false

// The following flags should be shared across preprocess and instrument.
const (
	WorkingDirEnv   = "OTEL_WORKING_DIRECTORY"
	DebugLogEnv     = "OTEL_DEBUG_TO_FILE"
	VerboseEnv      = "OTEL_VERBOSE"
	RuleJsonFileEnv = "OTEL_RULE_JSON_FILE"
)

// This is the version of the tool, which will be printed when the -version flag
// is passed. This value is specified by the build system.
var TheVersion = "1.0.0"

var TheName = "otelbuild"

func printVersion() {
	fmt.Printf("%s version %s\n", TheName, TheVersion)
}

func parseOptions() {
	// Parse flags from command-line arguments
	flag.BoolVar(&InToolexec, "in-toolexec", false,
		"Run in toolexec mode")
	flag.BoolVar(&DebugLog, "debuglog", false,
		"Print debug log to file")
	flag.BoolVar(&Verbose, "verbose", false,
		"Print verbose log")
	flag.BoolVar(&Debug, "debug", false,
		"Enable debug mode, leave temporary files for debugging")
	flag.BoolVar(&Restore, "restore", false,
		"Restore all instrumentations")
	flag.BoolVar(&PrintVersion, "version", false,
		"Print version")
	flag.StringVar(&RuleJsonFiles, "rule", "",
		"Rule file in json format. Multiple rules are separated by comma")
	flag.Parse()

	// Any non-flag command-line arguments will be treated as go build command
	GoBuildCmd = flag.Args()
}

func passOptions() error {
	if InInstrument() {
		// Inherit options from environment variables
		wd := os.Getenv(WorkingDirEnv)
		if len(wd) == 0 {
			return fmt.Errorf("cannot find working directory")
		}

		DebugLog, _ = strconv.ParseBool(os.Getenv(DebugLogEnv))
		Verbose, _ = strconv.ParseBool(os.Getenv(VerboseEnv))
		RuleJsonFiles = os.Getenv(RuleJsonFileEnv)
	} else {
		util.Assert(InPreprocess(), "why not otherwise")

		wd, err := filepath.Abs(TempBuildDir)
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		// Otherwise, set environment variables for further go toolexec build
		err = os.Setenv(WorkingDirEnv, wd)
		if err != nil {
			return fmt.Errorf("failed to set working directory: %w", err)
		}
		err = os.Setenv(DebugLogEnv, strconv.FormatBool(DebugLog))
		if err != nil {
			return fmt.Errorf("failed to set debug log flag: %w", err)
		}
		err = os.Setenv(VerboseEnv, strconv.FormatBool(Verbose))
		if err != nil {
			return fmt.Errorf("failed to set use rules flag: %w", err)
		}
	}
	return nil
}

func makeRuleAbs(file string) (string, error) {
	exist, err := util.PathExists(file)
	if err != nil {
		return "", fmt.Errorf("failed to check rule file: %w", err)
	}
	if !exist {
		return "", fmt.Errorf("rule file %s not found", file)
	}
	file, err = filepath.Abs(file)
	if err != nil {
		return "", fmt.Errorf("failed to get rule file: %w", err)
	}

	return file, nil
}

func makeRulesAbs() error {
	if RuleJsonFiles == "" {
		return nil
	}
	if strings.Contains(RuleJsonFiles, ",") {
		files := strings.Split(RuleJsonFiles, ",")
		for i, file := range files {
			f, err := makeRuleAbs(file)
			if err != nil {
				return fmt.Errorf("failed to set rule file: %w", err)
			}
			files[i] = f
		}
		RuleJsonFiles = strings.Join(files, ",")
	} else {
		f, err := makeRuleAbs(RuleJsonFiles)
		if err != nil {
			return fmt.Errorf("failed to set rule file: %w", err)
		}
		RuleJsonFiles = f
	}
	return nil
}

func checkOptions() error {
	if InPreprocess() {
		// Must be start with "go build", while "go" may refer to the full path
		if len(GoBuildCmd) < 2 ||
			!strings.Contains(GoBuildCmd[0], "go") ||
			GoBuildCmd[1] != "build" {
			return fmt.Errorf("usage: otelbuild go build")
		}
	}
	return nil
}

func InitOptions() (err error) {
	// Parse options from command-line arguments
	parseOptions()

	// Print version and exit early
	if PrintVersion {
		printVersion()
		os.Exit(0)
	}

	// Make sure all options are in sane
	err = checkOptions()
	if err != nil {
		return fmt.Errorf("failed to check options: %w", err)
	}

	// Pass options via environment variables for instrument phase
	err = passOptions()
	if err != nil {
		return fmt.Errorf("failed to inherit environment variables: %w", err)
	}

	if InPreprocess() {
		// Get absolute path of rule file, otherwise instrument will not
		// be able to find the rule file because it is running in different
		// working directory.
		err = makeRulesAbs()
		if err != nil {
			return fmt.Errorf("failed to set rule file: %w", err)
		}
		if err = os.Setenv(RuleJsonFileEnv, RuleJsonFiles); err != nil {
			return fmt.Errorf("failed to set rule file: %w", err)
		}
	}
	return nil
}
