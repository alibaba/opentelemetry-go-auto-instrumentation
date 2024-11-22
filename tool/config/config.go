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

package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
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

	// DebugLog true means debug log is enabled.
	DebugLog bool

	// Verbose true means print verbose log.
	Verbose bool

	// Debug true means debug mode.
	Debug bool

	// Restore true means restore all instrumentations.
	Restore bool

	disableDefaultRules bool
}

// This is the version of the tool, which will be printed when the -version flag
// is passed. This value is specified by the build system.
var ToolVersion = "1.0.0"

var GetToolName = func() string {
	// Get the path of the current executable
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable: %v", err)
		os.Exit(0)
	}
	return filepath.Base(ex)
}

var conf *BuildConfig

func GetConf() *BuildConfig {
	util.Assert(!shared.InConfigure(), "called in configure")
	util.Assert(conf != nil, "build config is not initialized")
	return conf
}

func (bc *BuildConfig) IsDisableDefaultRules() bool {
	return bc.disableDefaultRules
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
	if shared.InInstrument() {
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

func storeConfig(bc *BuildConfig) error {
	util.Assert(bc != nil, "build config is not initialized")
	util.Assert(shared.InConfigure(), "sanity check")

	file := shared.GetConfigureLogPath(shared.BuildConfFile)
	bs, err := json.Marshal(bc)
	if err != nil {
		return fmt.Errorf("failed to marshal build config: %w", err)
	}
	_, err = util.WriteFile(file, string(bs))
	if err != nil {
		return fmt.Errorf("failed to write build config: %w", err)
	}
	return nil
}

func loadConfig() (*BuildConfig, error) {
	util.Assert(conf == nil, "build config is already initialized")
	// If the build config file does not exist, return a default build config
	confFile := shared.GetConfigureLogPath(shared.BuildConfFile)
	exist, _ := util.PathExists(confFile)
	if !exist {
		return &BuildConfig{}, nil
	}
	// Load build config from json file
	file := shared.GetConfigureLogPath(shared.BuildConfFile)
	data, err := util.ReadFile(file)
	if err != nil {
		return &BuildConfig{},
			fmt.Errorf("failed to read build config: %w", err)
	}
	bc := &BuildConfig{}
	err = json.Unmarshal([]byte(data), bc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal build config: %w", err)
	}
	return bc, nil
}

func InitConfig() (err error) {
	// Load build config from json file
	conf, err = loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load build config: %w", err)
	}

	err = conf.parseRuleFiles()
	if err != nil {
		return fmt.Errorf("failed to parse rule files: %w", err)
	}

	if conf.DebugLog {
		// Redirect log to debug log if required
		debugLogPath := shared.GetLogPath(shared.DebugLogFile)
		debugLog, _ := os.OpenFile(debugLogPath, os.O_WRONLY|os.O_APPEND, 0777)
		if debugLog != nil {
			log.SetOutput(debugLog)
		}
	}
	return nil
}

func PrintVersion() error {
	fmt.Printf("%s version %s", GetToolName(), ToolVersion)
	os.Exit(0)
	return nil
}

func Configure() error {
	shared.GuaranteeInConfigure()

	// Parse command line flags to get build config
	bc, err := loadConfig()
	if err != nil {
		bc = &BuildConfig{}
	}
	flag.BoolVar(&bc.DebugLog, "debuglog", bc.DebugLog,
		"Print debug log to file")
	flag.BoolVar(&bc.Verbose, "verbose", bc.Verbose,
		"Print verbose log")
	flag.BoolVar(&bc.Debug, "debug", bc.Debug,
		"Enable debug mode, leave temporary files for debugging")
	flag.BoolVar(&bc.Restore, "restore", bc.Restore,
		"Restore all instrumentations")
	flag.StringVar(&bc.RuleJsonFiles, "rule", bc.RuleJsonFiles,
		"Rule file in json format. Multiple rules are separated by comma")
	flag.CommandLine.Parse(os.Args[2:])

	fmt.Printf("Configured in %s",
		shared.GetConfigureLogPath(shared.BuildConfFile))

	// Store build config for future phases
	err = storeConfig(bc)
	if err != nil {
		return fmt.Errorf("failed to store build config: %w", err)
	}
	return nil
}
