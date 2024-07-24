package tool

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/instrument"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/preprocess"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
)

func initLogs(names ...string) error {
	for _, name := range names {
		path := filepath.Join(shared.TempBuildDir, name)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		if shared.DebugLog {
			logPath := filepath.Join(path, shared.DebugLogFile)
			_, err = os.Create(logPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setupLogs() {
	if shared.InInstrument() {
		log.SetPrefix("[" + shared.TInstrument + "] ")
	} else {
		log.SetPrefix("[" + shared.TPreprocess + "] ")
	}
	if shared.DebugLog {
		// Redirect log to debug log if required
		debugLogPath := shared.GetLogPath(shared.DebugLogFile)
		debugLog, _ := os.OpenFile(debugLogPath, os.O_WRONLY|os.O_APPEND, 0777)
		if debugLog != nil {
			log.SetOutput(debugLog)
		}
	}
}

func initEnv() (err error) {
	if shared.InInstrument() {
		setupLogs()
	} else {
		err = os.MkdirAll(shared.TempBuildDir, 0777)
		if err != nil {
			return fmt.Errorf("failed to make working directory: %w", err)
		}

		// @@ Init here to avoid permission issue
		err = initLogs(shared.TPreprocess, shared.TInstrument)
		if err != nil {
			return fmt.Errorf("failed to init logs: %w", err)
		}

		setupLogs()
	}
	err = resource.InitRules()
	if err != nil {
		return fmt.Errorf("failed to init rules: %w", err)
	}

	// Disable all instrumentation rules and rebuild the whole project to restore
	// all instrumentation actions, this also reverts the modification on Golang
	// runtime package.
	if shared.Restore {
		shared.DisableRules = "*"
	}
	return nil
}

func Run() (err error) {
	// Where our story begins
	err = initEnv()
	if err != nil {
		return fmt.Errorf("failed to init context: %w", err)
	}

	if shared.InPreprocess() {
		return preprocess.Preprocess()
	} else {
		return instrument.Instrument()
	}
}
