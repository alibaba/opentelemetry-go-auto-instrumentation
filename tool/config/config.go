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
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	EnvPrefix     = "OTELTOOL_"
	BuildConfFile = "conf.json"
)

type BuildConfig struct {
	// RuleJsonFiles is the name of the rule file. It is used to tell instrument
	// tool where to find the instrument rules. Multiple rules are separated by
	// comma. e.g. -rule=rule1.json,rule2.json. By default, new rules are appended
	// to default rules, i.e. -rule=rule1.json,rule2.json is exactly equivalent to
	// -rule=default.json,rule1.json,rule2.json. But if you do want to disable
	// default rules, you can configure -disabledefault flag in advance.
	RuleJsonFiles string

	// Log specifies the log file path. If not set, log will be saved to file.
	Log string

	// Verbose true means print verbose log.
	Verbose bool

	// Debug true means debug mode.
	Debug bool

	// Restore true means restore all instrumentations.
	Restore bool

	// DisableDefault true means disable default rules.
	DisableDefault bool
}

// @@This value is specified by the build system.
// This is the version of the tool, which will be printed when the -version flag
// is passed.
var ToolVersion = "1.0.0"

// @@This value is specified by the build system.
// It specifies the source path of the tool, which will be used to find the rules
var BuildPath = ""

// @@This value is specified by the build system.
// It specifies the version of the pkg module, whose rules resides in it.
// If the value is "latest", it means the latest version of the pkg module will
// be used. If the value is a specific version, it means the specific version
// of the pkg module will be used.
// We added this flag because we want each release of the otel tool to precisely
// bind to a specific version of the pkg module. Without this flag, every version
// of the otel tool would pull the latest pkg modules (i.e., pkg/rules), which
// is not our intention.
var UsedPkg = "latest"

var conf *BuildConfig

func GetConf() *BuildConfig {
	util.Assert(conf != nil, "build config is not initialized")
	return conf
}

func (bc *BuildConfig) IsDisableDefault() bool {
	return bc.DisableDefault
}

func (bc *BuildConfig) makeRuleAbs(file string) (string, error) {
	if util.PathNotExists(file) {
		return "", errc.New(errc.ErrNotExist, file)
	}
	file, err := filepath.Abs(file)
	if err != nil {
		return "", errc.New(errc.ErrAbsPath, err.Error())
	}
	return file, nil
}

func (bc *BuildConfig) parseRuleFiles() error {
	if util.InInstrument() {
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
				return err
			}
			files[i] = f
		}
		bc.RuleJsonFiles = strings.Join(files, ",")
	} else {
		f, err := bc.makeRuleAbs(bc.RuleJsonFiles)
		if err != nil {
			return err
		}
		bc.RuleJsonFiles = f
	}
	return nil
}

func getConfPath(name string) string {
	return util.GetTempBuildDirWith(name)
}

func storeConfig(bc *BuildConfig) error {
	util.Assert(bc != nil, "build config is not initialized")

	file := getConfPath(BuildConfFile)
	bs, err := json.Marshal(bc)
	if err != nil {
		return errc.New(errc.ErrInvalidJSON, err.Error())
	}
	_, err = util.WriteFile(file, string(bs))
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() (*BuildConfig, error) {
	util.Assert(conf == nil, "build config is already initialized")
	// If the build config file does not exist, return a default build config
	confFile := getConfPath(BuildConfFile)
	if util.PathNotExists(confFile) {
		return &BuildConfig{}, nil
	}
	// Load build config from json file
	file := getConfPath(BuildConfFile)
	data, err := util.ReadFile(file)
	if err != nil {
		return &BuildConfig{}, err
	}
	bc := &BuildConfig{}
	err = json.Unmarshal([]byte(data), bc)
	if err != nil {
		return nil, errc.New(errc.ErrInvalidJSON, err.Error())
	}
	return bc, nil
}

func toUpperSnakeCase(input string) string {
	var result []rune

	for i, char := range input {
		if unicode.IsUpper(char) {
			if i != 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToUpper(char))
		} else {
			result = append(result, unicode.ToUpper(char))
		}
	}

	return string(result)
}
func loadConfigFromEnv(conf *BuildConfig) {
	// Environment variables are able to overwrite the config items even if the
	// config file sets them. The environment variable name is the upper snake
	// case of the config item name, prefixed with "OTELTOOL_". For example, the
	// environment variable for "Log" is "OTELTOOL_LOG".
	typ := reflect.TypeOf(*conf)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		envKey := fmt.Sprintf("%s%s", EnvPrefix, toUpperSnakeCase(field.Name))
		envVal := os.Getenv(envKey)
		if envVal != "" {
			if util.InPreprocess() {
				util.Log("Overwrite config %s with environment variable %s",
					field.Name, envKey)
			}
			v := reflect.ValueOf(conf).Elem()
			f := v.FieldByName(field.Name)
			switch f.Kind() {
			case reflect.Bool:
				f.SetBool(envVal == "true")
			case reflect.String:
				f.SetString(envVal)
			default:
				util.ShouldNotReachHere()
			}
		}
	}
}

func InitConfig() (err error) {
	// Load build config from json file
	conf, err = loadConfig()
	if err != nil {
		return err
	}
	loadConfigFromEnv(conf)

	err = conf.parseRuleFiles()
	if err != nil {
		return err
	}

	mode := os.O_WRONLY | os.O_APPEND
	if util.InPreprocess() {
		// We always create log file in preprocess phase, but in further
		// instrument phase, we append log content to the existing file.
		mode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	if conf.Log == "" {
		// Redirect log to file if flag is not set
		debugLogPath := util.GetTempBuildDirWith(util.DebugLogFile)
		debugLog, _ := os.OpenFile(debugLogPath, mode, 0777)
		if debugLog != nil {
			util.SetLogger(debugLog)
		}
	} else {
		// Otherwise, log to the specified file
		logFile, err := os.OpenFile(conf.Log, mode, 0777)
		if err != nil {
			return errc.New(errc.ErrOpenFile, err.Error())
		}
		util.SetLogger(logFile)
	}
	return nil
}

func PrintVersion() error {
	name, err := util.GetToolName()
	if err != nil {
		return err
	}
	fmt.Printf("%s version %s\n", name, ToolVersion)
	return nil
}

func Configure() error {
	// Parse command line flags to get build config
	bc, err := loadConfig()
	if err != nil {
		bc = &BuildConfig{}
	}
	flag.StringVar(&bc.Log, "log", bc.Log,
		"Log file path. If not set, log will be saved to file.")
	flag.BoolVar(&bc.Verbose, "verbose", bc.Verbose,
		"Print verbose log")
	flag.BoolVar(&bc.Debug, "debug", bc.Debug,
		"Enable debug mode, leave temporary files for debugging")
	flag.BoolVar(&bc.Restore, "restore", bc.Restore,
		"Restore all instrumentations")
	flag.StringVar(&bc.RuleJsonFiles, "rule", bc.RuleJsonFiles,
		"Use custom.json rules. Multiple rules are separated by comma.")
	flag.BoolVar(&bc.DisableDefault, "disabledefault", bc.DisableDefault,
		"Disable default rules")
	flag.CommandLine.Parse(os.Args[2:])

	util.Log("Configured in %s", getConfPath(BuildConfFile))

	// Store build config for future phases
	err = storeConfig(bc)
	if err != nil {
		return err
	}
	return nil
}
