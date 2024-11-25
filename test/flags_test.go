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

func TestFlags(t *testing.T) {
	const AppName = "flags"
	UseApp(AppName)

	RunGoBuildFallible(t, "go", "build", "-thisisnotvalid")
	ExpectPreprocessContains(t, shared.DebugLogFile, "failed to")

	RunVersion(t)
	ExpectStdoutContains(t, "version")

	RunGoBuildFallible(t, "go", "build", "notevenaflag")
	ExpectPreprocessContains(t, shared.DebugLogFile, "failed to")

	RunSet(t, "-verbose")
	RunGoBuild(t, "go", "build", `-ldflags=-X main.Placeholder=replaced`)
	_, stderr := RunApp(t, AppName)
	ExpectContains(t, stderr, "placeholder:replaced")

	RunGoBuild(t, "go")
	RunGoBuild(t)
	RunGoBuild(t, "")
}
