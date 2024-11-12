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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	TempBuildDir    = ".otel-build"
	VendorDir       = "vendor"
	BuildModeVendor = "-mod=vendor"
	BuildModeMod    = "-mod=mod"
	BuildConfFile   = "build_conf.json"
)

type BuildConfig struct {
	// RuleJsonFiles is the name of the rule file. It is used to tell instrument
	// tool where to find the instrument rules. Multiple rules are separated by
	// comma. e.g. -rule=rule1.json,rule2.json. By default, new rules are appended
	// to default rules, i.e. -rule=rule1.json,rule2.json is exactly equivalent to
	// -rule=default.json,rule1.json,rule2.json. But if you do want to replace the
	// default rules, you can add a "+" prefix to the rule file name, e.g.
	// -rule=+rule1.json,rule2.json. In this case, the default rules will be replaced
	// by rule1.json, and then rule2.json will be appended to the rules.
	RuleJsonFiles string

	// InToolexec true means this tool is being invoked in the go build process.
	// This flag **SHOULD NOT** be set manually by users.
	InToolexec bool

	// DebugLog true means debug log is enabled.
	DebugLog bool

	// Verbose true means print verbose log.
	Verbose bool

	// Debug true means debug mode.
	Debug bool

	// Restore true means restore all instrumentations.
	Restore bool

	// GoBuildCmd is equivalent to go build command
	GoBuildCmd []string

	// PrintVersion true means print version.
	PrintVersion bool

	disableDefaultRules bool
}

// This is the version of the tool, which will be printed when the -version flag
// is passed. This value is specified by the build system.
var TheVersion = "1.0.0"

var TheName = "otelbuild"

var conf *BuildConfig

func printVersion() {
	util.Assert(conf != nil, "build config is not initialized")

	if conf.PrintVersion {
		fmt.Printf("%s version %s\n", TheName, TheVersion)
		os.Exit(0)
	}
}

func GetConf() *BuildConfig {
	util.Assert(conf != nil, "build config is not initialized")
	return conf
}

func (bc *BuildConfig) IsDisableDefaultRules() bool {
	return bc.disableDefaultRules
}

func (bc *BuildConfig) setBuildMode() error {
	// We can't brutely always add -mod=mod here because -mod may only be set to
	// readonly or vendor when in workspace mode. We need to check if provided
	// -mod is vendor or vendor directory exists, then we set -mod=mod mode.
	// For all other cases, we just leave it as is.

	// Check if -mod=vendor is set, replace it with -mod=mod
	const buildModeVendor = "-mod=vendor"
	const buildModePrefix = "-mod"
	for i, arg := range bc.GoBuildCmd {
		// -mod=vendor?
		if arg == buildModeVendor {
			bc.GoBuildCmd[i] = BuildModeMod
			return nil
		}
		// -mod vendor?
		if arg == buildModePrefix {
			if len(bc.GoBuildCmd) > i+1 {
				arg1 := bc.GoBuildCmd[i+1]
				if arg1 == "vendor" {
					bc.GoBuildCmd[i+1] = "mod"
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
		bc.GoBuildCmd = append([]string{BuildModeMod}, bc.GoBuildCmd...)
	}
	return nil
}

func (bc *BuildConfig) makeRuleAbs(file string) (string, error) {
	// Check if rule json file has a "+" prefix, which means to replace the
	// default rules, i.e. whether to keep the default rules.
	if strings.HasPrefix(file, "+") {
		bc.disableDefaultRules = true
		file = file[1:]
	}
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

func (bc *BuildConfig) parseRuleFiles() error {
	if InInstrument() {
		return nil
	}
	// Get absolute path of rule file, otherwise instrument will not
	// be able to find the rule file because it is running in different
	// working directory.
	if bc.RuleJsonFiles == "" {
		return nil
	}
	if strings.Contains(bc.RuleJsonFiles, ",") {
		files := strings.Split(bc.RuleJsonFiles, ",")
		for i, file := range files {
			f, err := bc.makeRuleAbs(file)
			if err != nil {
				return fmt.Errorf("failed to set rule file: %w", err)
			}
			files[i] = f
		}
		bc.RuleJsonFiles = strings.Join(files, ",")
	} else {
		f, err := bc.makeRuleAbs(bc.RuleJsonFiles)
		if err != nil {
			return fmt.Errorf("failed to set rule file: %w", err)
		}
		bc.RuleJsonFiles = f
	}
	return nil
}

func storeBuildConfig() error {
	util.Assert(conf != nil, "build config is not initialized")
	util.Assert(InPreprocess(), "sanity check")

	file := GetPreprocessLogPath(BuildConfFile)
	bs, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to marshal build config: %w", err)
	}
	_, err = util.WriteFile(file, string(bs))
	if err != nil {
		return fmt.Errorf("failed to write build config: %w", err)
	}
	return nil
}

func loadBuildConfig() (*BuildConfig, error) {
	util.Assert(conf == nil, "build config is already initialized")

	// Early initilaization for phase identification
	bc := &BuildConfig{}
	flag.BoolVar(&bc.InToolexec, "in-toolexec", false,
		"Run in toolexec mode")

	// In Preprocess phase, we parse build config from command-line arguments.
	if !bc.InToolexec {
		flag.BoolVar(&bc.DebugLog, "debuglog", false,
			"Print debug log to file")
		flag.BoolVar(&bc.Verbose, "verbose", false,
			"Print verbose log")
		flag.BoolVar(&bc.Debug, "debug", false,
			"Enable debug mode, leave temporary files for debugging")
		flag.BoolVar(&bc.Restore, "restore", false,
			"Restore all instrumentations")
		flag.BoolVar(&bc.PrintVersion, "version", false,
			"Print version")
		flag.StringVar(&bc.RuleJsonFiles, "rule", "",
			"Rule file in json format. Multiple rules are separated by comma")
		flag.Parse()

		// Any non-flag command-line arguments will be treated as build
		// arguments and transparently passed to the go build command.
		bc.GoBuildCmd = flag.Args()

		// At this point, the build config is fully initialized and ready to use.
		return bc, nil
	} else {
		// In Instrument phase, we should not parse the flags, instead we load
		// the config from json file.
		file := GetPreprocessLogPath(BuildConfFile)
		data, err := util.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read build config: %w", err)
		}
		err = json.Unmarshal([]byte(data), bc)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal build config: %w", err)
		}
		return bc, nil
	}
}

func initLogs(names ...string) error {
	for _, name := range names {
		path := filepath.Join(TempBuildDir, name)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		if GetConf().DebugLog {
			logPath := filepath.Join(path, DebugLogFile)
			_, err = os.Create(logPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setupLogs() {
	if InInstrument() {
		log.SetPrefix("[" + TInstrument + "] ")
	} else {
		log.SetPrefix("[" + TPreprocess + "] ")
	}
	if GetConf().DebugLog {
		// Redirect log to debug log if required
		debugLogPath := GetLogPath(DebugLogFile)
		debugLog, _ := os.OpenFile(debugLogPath, os.O_WRONLY|os.O_APPEND, 0777)
		if debugLog != nil {
			log.SetOutput(debugLog)
		}
	}
}

func initTempDir() error {
	// Create temp build directory in preprocess phase, this should be
	// done in very early stage, as every further operation will likely
	// depend on this directory.
	if InPreprocess() {
		// Clean up temp build directory if exists, otherwise create it
		_, err := os.Stat(TempBuildDir)
		if err == nil {
			err = os.RemoveAll(TempBuildDir)
			if err != nil {
				return fmt.Errorf("failed to remove working directory: %w", err)
			}
		}
		err = os.MkdirAll(TempBuildDir, 0777)
		if err != nil {
			return fmt.Errorf("failed to make working directory: %w", err)
		}
		// @@ Init here to avoid permission issue
		err = initLogs(TPreprocess, TInstrument)
		if err != nil {
			return fmt.Errorf("failed to init logs: %w", err)
		}

	}

	setupLogs()
	return nil
}

func InitConfig() (err error) {
	// Load build config from either command-line arguments or json file
	conf, err = loadBuildConfig()
	if err != nil {
		return fmt.Errorf("failed to load build config: %w", err)
	}

	// Init temp build directory in very early stage during preprocess phase
	err = initTempDir()
	if err != nil {
		return fmt.Errorf("failed to init temp dir: %w", err)
	}

	// Print version and exit early
	printVersion()

	err = conf.parseRuleFiles()
	if err != nil {
		return fmt.Errorf("failed to parse rule files: %w", err)
	}

	err = conf.setBuildMode()
	if err != nil {
		return fmt.Errorf("failed to set build mode: %w", err)
	}

	// Store build config to json file for instrument phase
	if InPreprocess() {
		err = storeBuildConfig()
		if err != nil {
			return fmt.Errorf("failed to store build config: %w", err)
		}
	}
	return nil
}
