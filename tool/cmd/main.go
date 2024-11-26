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
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/instrument"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/preprocess"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
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
	usage = strings.ReplaceAll(usage, "{}", util.GetToolName())
	fmt.Print(usage)
}

func initLog() error {
	name := util.PPreprocess
	path := shared.GetTempBuildDirWith(name)
	logPath := filepath.Join(path, shared.DebugLogFile)
	_, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
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
	if exist, _ := util.PathExists(shared.TempBuildDir); !exist {
		err := os.MkdirAll(shared.TempBuildDir, 0777)
		if err != nil {
			return fmt.Errorf("failed to make working directory: %w", err)
		}
	}
	// Make sub-directory of temp build directory for each phase. Specifaclly,
	// we always recreate the preprocess and instrument directories, but only
	// create the configure directory if it does not exist. This is because
	// the configure directory can be used across multiple runs.
	exist, _ := util.PathExists(shared.GetTempBuildDirWith(util.PConfigure))
	if !exist {
		err := os.MkdirAll(shared.GetTempBuildDirWith(util.PConfigure), 0777)
		if err != nil {
			return fmt.Errorf("failed to make log directory: %w", err)
		}
	}
	for _, subdir := range []string{util.PPreprocess, util.PInstrument} {
		exist, _ = util.PathExists(shared.GetTempBuildDirWith(subdir))
		if exist {
			err := os.RemoveAll(shared.GetTempBuildDirWith(subdir))
			if err != nil {
				return fmt.Errorf("failed to remove directory: %w", err)
			}
		}
		err := os.MkdirAll(shared.GetTempBuildDirWith(subdir), 0777)
		if err != nil {
			return fmt.Errorf("failed to make log directory: %w", err)
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
		return fmt.Errorf("failed to init temp dir: %w", err)
	}

	// Create log files under temp build directory
	if util.InPreprocess() {
		err := initLog()
		if err != nil {
			return fmt.Errorf("failed to init logs: %w", err)
		}
	}

	// Prepare shared configuration
	if util.InPreprocess() || util.InInstrument() {
		err = config.InitConfig()
		if err != nil {
			return fmt.Errorf("failed to init config: %w", err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	err := initEnv()
	if err != nil {
		util.LogFatal("failed to init env: %v", err)
		os.Exit(1)
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
		util.LogFatal("failed to run command %s: %v", subcmd, err)
		os.Exit(1)
	}
}
