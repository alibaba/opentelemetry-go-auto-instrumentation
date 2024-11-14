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

package tool

import (
	"log"
	"os"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/instrument"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/preprocess"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
)

func setupLogs() {
	if shared.InInstrument() {
		log.SetPrefix("[" + shared.TInstrument + "] ")
	} else {
		log.SetPrefix("[" + shared.TPreprocess + "] ")
	}
	if shared.GetBuildConfig().DebugLog {
		// Redirect log to debug log if required
		debugLogPath := shared.GetLogPath(shared.DebugLogFile)
		debugLog, _ := os.OpenFile(debugLogPath, os.O_WRONLY|os.O_APPEND, 0777)
		if debugLog != nil {
			log.SetOutput(debugLog)
		}
	}
}

func Build() (err error) {
	// Where our story begins
	setupLogs()
	if shared.InPreprocess() {
		return preprocess.Preprocess()
	} else {
		return instrument.Instrument()
	}
}
