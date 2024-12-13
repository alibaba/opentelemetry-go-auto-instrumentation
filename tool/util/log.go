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
)

var logWriter *os.File = os.Stdout

var Guarantee = Assert // More meaningful name:)

func SetLogTo(w *os.File) {
	logWriter = w
}

func Log(format string, args ...interface{}) {
	template := "[" + GetRunPhase().String() + "] " + format + "\n"
	fmt.Fprintf(logWriter, template, args...)
}

func LogFatal(format string, args ...interface{}) {
	// Log errors to debug file
	Log(format, args...)
	// And print to stderr then, in red color
	template := "Build error:\033[31m\n" + format + "\033[0m\n"
	fmt.Fprintf(os.Stderr, template,
		args...)
	fmt.Fprintf(os.Stderr, "See build log %s for details.\n",
		logWriter.Name())
	os.Exit(1)
}
