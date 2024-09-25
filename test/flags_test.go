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

package test

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"testing"
)

func TestFlags(t *testing.T) {
	const AppName = "flags"
	UseApp(AppName)

	RunInstrument(t, "-debuglog", "-disablerules=")
	ExpectInstrumentContains(t, shared.DebugLogFile, "fmt@")

	RunInstrument(t, "-debuglog", "-disablerules=fmt,net/http")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "fmt@")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=*")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "fmt@")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=testrule")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "fmt@")

	RunInstrumentFallible(t, "-debuglog", "--", "-thisisnotvalid")
	ExpectPreprocessContains(t, shared.DebugLogFile, "failed to")

	RunInstrument(t, "-version")
	ExpectStdoutContains(t, "version")

	RunInstrumentFallible(t, "-debuglog", "--", "notevenaflag")
	ExpectPreprocessContains(t, shared.DebugLogFile, "failed to")

	RunInstrument(t, "-debuglog", "-verbose",
		"--",
		`-ldflags=-X main.Placeholder=replaced`)
	_, stderr := RunApp(t, AppName)
	ExpectContains(t, stderr, "placeholder:replaced")
}
