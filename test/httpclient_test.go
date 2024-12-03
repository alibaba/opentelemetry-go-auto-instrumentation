// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"
)

const HttpclientAppName = "httpclient"

func TestRunHttpclient(t *testing.T) {
	UseApp(HttpclientAppName)

	RunSet(t, UseTestRules("test_nethttp.json"), "-verbose")
	RunGoBuild(t, "go", "build")
	_, stderr := RunApp(t, HttpclientAppName)
	ExpectContains(t, stderr, "Client.Do()")                // println writes to stderr
	ExpectContains(t, stderr, "failed to exec onExit hook") // intentional panic
	ExpectContains(t, stderr, "NewRequest()")
	ExpectContains(t, stderr, "NewRequest1()")
	ExpectContains(t, stderr, "NewRequestWithContext()")
	ExpectContains(t, stderr, "MaxBytesError()")
	ExpectContains(t, stderr, "debug.Stack()") // during recover()
	ExpectContains(t, stderr, "4008208820")
	ExpectContains(t, stderr, "Prince of Qin Smashing the Battle line")
}
