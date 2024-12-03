// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
