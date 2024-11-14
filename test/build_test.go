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
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
)

func TestBuildProject(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)
	RunInstrument(t, "--", "-o", "default", "cmd/foo.go")
	RunInstrument(t, "--", "cmd/foo.go")
	RunInstrument(t, "--", "cmd/foo.go", "cmd/bar.go")
	RunInstrument(t, "--", "cmd")
}

func TestBuildProject2(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "--", ".")
	RunInstrument(t, "--", "")
}

func TestBuildProject3(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "--", "m1")
	RunInstrumentFallible(t, "--", "m2") // not used in go.work
}

func TestBuildProject4(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "-rule=+../../pkg/data/default.json", "--", "m1")
	RunInstrumentFallible(t, "-rule=../../pkg/data/default.json", "--", "m1") // duplicated default rules
	RunInstrumentFallible(t, "-rule=../../pkg/data/default", "--", "m1")
	RunInstrument(t, "-rule=../../pkg/data/test_error.json,../../pkg/data/test_fmt.json", "--", "m1")
	RunInstrument(t, "-rule=../../pkg/data/test_error.json,+../../pkg/data/test_fmt.json", "--", "m1")
	RunInstrument(t, "-rule=+../../pkg/data/default.json,+../../pkg/data/test_fmt.json", "--", "m1")
}

func TestBuildProject5(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "-debuglog", "-rule=../../pkg/data/test_fmt.json",
		"-verbose", "--", "m1")
	// both test_fmt.json and default.json rules should be available
	// because we always append new -rule to the default.json by default
	// (unless we use -rule=+... syntax) to explicitly disable default rules.
	ExpectPreprocessContains(t, shared.DebugLogFile, "fmt")
	ExpectPreprocessContains(t, shared.DebugLogFile, "database/sql")
}

func TestBuildProject6(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "-debuglog", "-rule=+../../pkg/data/test_fmt.json",
		"-verbose", "--", "m1")
	// both test_fmt.json and default.json rules should be available
	// because we always append new -rule to the default.json by default
	// (unless we use -rule=+... syntax) to explicitly disable default rules.
	ExpectPreprocessContains(t, shared.DebugLogFile, "fmt")
	ExpectPreprocessNotContains(t, shared.DebugLogFile, "database/sql")
}
