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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/instrument"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/preprocess"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	SubcommandSet     = "set"
	SubcommandGo      = "go"
	SubcommandVersion = "version"
	SubcommandRemix   = "remix"
)

var usage = `Usage: {} <command> [args]
Example:
	{} go build
	{} go build main.go
	{} version
	{} set -verbose -rule=custom.json

Command:
	version    print the version
	set        set the configuration
	go         build the Go application
`

func printUsage() {
	name, _ := util.GetToolName()
	usage = strings.ReplaceAll(usage, "{}", name)
	fmt.Print(usage)
}

func initLog() error {
	name := util.PPreprocess
	path := util.GetTempBuildDirWith(name)
	logPath := filepath.Join(path, util.DebugLogFile)
	_, err := os.Create(logPath)
	if err != nil {
		return errc.New(errc.ErrCreateFile, err.Error())
	}
	return nil
}

func initTempDir() error {
	// All temp directories are prepared before, instrument phase should not
	// create any new directories.
	if util.GetRunPhase() == util.PInstrument {
		return nil
	}

	// Make temp build directory if not exists
	if util.PathNotExists(util.TempBuildDir) {
		err := os.MkdirAll(util.TempBuildDir, 0777)
		if err != nil {
			return errc.New(errc.ErrMkdirAll, err.Error())
		}
	}
	// Make sub-directory of temp build directory for each phase. Specifaclly,
	// we always recreate the preprocess and instrument directories, but only
	// create the configure directory if it does not exist. This is because
	// the configure directory can be used across multiple runs.
	if util.PathNotExists(util.GetTempBuildDirWith(util.PConfigure)) {
		err := os.MkdirAll(util.GetTempBuildDirWith(util.PConfigure), 0777)
		if err != nil {
			return errc.New(errc.ErrMkdirAll, err.Error())
		}
	}
	for _, subdir := range []string{util.PPreprocess, util.PInstrument} {
		if util.PathExists(util.GetTempBuildDirWith(subdir)) {
			err := os.RemoveAll(util.GetTempBuildDirWith(subdir))
			if err != nil {
				return errc.New(errc.ErrRemoveAll, err.Error())
			}
		}
		err := os.MkdirAll(util.GetTempBuildDirWith(subdir), 0777)
		if err != nil {
			return errc.New(errc.ErrMkdirAll, err.Error())
		}
	}

	return nil
}

func initEnv() error {
	util.Assert(len(os.Args) >= 2, "no command specified")

	// Determine the run phase
	switch {
	case os.Args[1] == SubcommandSet:
		// otel set?
		util.SetRunPhase(util.PConfigure)
	case strings.HasSuffix(os.Args[1], SubcommandGo):
		// otel go build?
		util.SetRunPhase(util.PPreprocess)
	case os.Args[1] == SubcommandRemix:
		// otel remix?
		util.SetRunPhase(util.PInstrument)
	default:
		// do nothing
	}

	// Create temp build directory
	err := initTempDir()
	if err != nil {
		return err
	}

	// Create log files under temp build directory
	if util.InPreprocess() {
		err := initLog()
		if err != nil {
			return err
		}
	}

	// Prepare shared configuration
	if util.InPreprocess() || util.InInstrument() {
		err = config.InitConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

func fatal(err error) {
	message := "===== Environments =====\n"
	message += fmt.Sprintf("%-11s: %s\n", "Command", strings.Join(os.Args, " "))
	message += fmt.Sprintf("%-11s: %s\n", "ErrorLog", util.GetLoggerPath())
	message += fmt.Sprintf("%-11s: %s\n", "WorkDir", os.Getenv("PWD"))
	message += fmt.Sprintf("%-11s: %s, %s, %s\n", "Toolchain",
		runtime.GOOS+"/"+runtime.GOARCH,
		runtime.Version(), config.ToolVersion)
	message += "===== Fatal Error ======\n"
	message += err.Error()
	util.LogFatal("\033[31m%s\033[0m", message) // log in red color
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	err := initEnv()
	if err != nil {
		fatal(err)
	}

	subcmd := os.Args[1]
	switch subcmd {
	case SubcommandVersion:
		err = config.PrintVersion()
	case SubcommandSet:
		err = config.Configure()
	case SubcommandGo:
		err = preprocess.Preprocess()
	case SubcommandRemix:
		err = instrument.Instrument()
	default:
		printUsage()
	}
	if err != nil {
		if subcmd != SubcommandRemix {
			fatal(err)
		} else {
			// If error occurs in remix phase, we dont want to decoret the error
			// message with the environments, just print the error message, the
			// caller(preprocess) phase will decorate instead.
			util.LogFatal(err.Error())
		}
	}
}
