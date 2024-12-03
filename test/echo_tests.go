// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

const echo_dependency_name = "github.com/labstack/echo/v4"
const echo_module_name = "echo"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("echo-basic-test", echo_module_name, "v4.0.0", "", "1.18", "", TestBasicEcho),
		NewGeneralTestCase("echo-middleware-test", echo_module_name, "v4.0.0", "", "1.18", "", TestEchoMiddleware),
		NewGeneralTestCase("echo-pattern-test", echo_module_name, "v4.0.0", "", "1.18", "", TestEchoPattern),
		NewMuzzleTestCase("echo-muzzle-test", echo_dependency_name, echo_module_name, "v4.0.0", "v4.9.1", "1.18", "", []string{"go", "build", "test_echo_basic.go"}),
		NewMuzzleTestCase("echo-muzzle-test", echo_dependency_name, echo_module_name, "v4.10.0", "", "1.18", "", []string{"go", "build", "test_echo_middleware.go"}),
		NewLatestDepthTestCase("echo-latestdepth-test", echo_dependency_name, echo_module_name, "v4.10.0", "", "1.18", "", TestBasicEcho),
	)
}

func TestBasicEcho(t *testing.T, env ...string) {
	UseApp("echo/v4.0.0")
	RunGoBuild(t, "go", "build", "test_echo_basic.go")
	RunApp(t, "test_echo_basic", env...)
}

func TestEchoPattern(t *testing.T, env ...string) {
	UseApp("echo/v4.0.0")
	RunGoBuild(t, "go", "build", "test_echo_pattern.go")
	RunApp(t, "test_echo_pattern", env...)
}

func TestEchoMiddleware(t *testing.T, env ...string) {
	UseApp("echo/v4.10.0")
	RunGoBuild(t, "go", "build", "test_echo_middleware.go")
	RunApp(t, "test_echo_middleware", env...)
}
