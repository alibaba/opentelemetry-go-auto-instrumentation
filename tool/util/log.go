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

package util

import (
	"fmt"
	"os"
	"sync"
)

var logWriter *os.File = os.Stdout
var logMutex sync.Mutex

var Guarantee = Assert // More meaningful name:)

func SetLogTo(w *os.File) {
	logWriter = w
}

func Log(format string, args ...interface{}) {
	template := "[" + GetRunPhase().String() + "] " + format + "\n"
	logMutex.Lock()
	fmt.Fprintf(logWriter, template, args...)
	logMutex.Unlock()
}

func LogFatal(format string, args ...interface{}) {
	// Log errors to debug file
	Log(format, args...)
	// And print to stderr then, in red color
	if InPreprocess() {
		// Print more information for preprocess
		details := map[string]string{}
		details["Command"] = os.Args[0]
		details["ErrorLog"] = logWriter.Name()
		details["WorkDir"] = os.Getenv("PWD")

		fmt.Fprintf(os.Stderr, "%-10s: ",
			"BuildError")
		fmt.Fprintf(os.Stderr, "\033[31m"+format+"\033[0m", args...)
		for name, value := range details {
			fmt.Fprintf(os.Stderr, "%-10s: %s\n", name, value)
		}
	}
	os.Exit(1)
}
