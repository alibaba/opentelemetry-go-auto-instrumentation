// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

const hertz_dependency_name = "github.com/cloudwego/hertz"
const hertz_module_name = "hertz"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("hertz-090-basic-test", hertz_module_name, "v0.9.0", "", "1.18", "1.22", TestBasicHertz),
		NewGeneralTestCase("hertz-090-basic-test-with-hook", hertz_module_name, "v0.9.0", "", "1.18", "1.22", TestBasicHertzWithHook),
		NewGeneralTestCase("hertz-090-basic-test-with-exception", hertz_module_name, "v0.9.0", "", "1.18", "1.22", TestBasicHertzWithException),
		NewGeneralTestCase("hertz-090-basic-test-with-regex", hertz_module_name, "v0.9.0", "", "1.18", "1.22", TestBasicHertzWithRegex),
		NewLatestDepthTestCase("hertz-090-basic-test-latestdepth", hertz_dependency_name, hertz_module_name, "v0.9.0", "", "1.18", "", TestBasicHertz),
		NewMuzzleTestCase("hertz-090-basic-muzzle", hertz_dependency_name, hertz_module_name, "v0.9.0", "v0.9.1", "1.18", "1.22.9", []string{"go", "build", "test_hertz_basic.go", "basic_func.go"}),
		NewMuzzleTestCase("hertz-090-basic-muzzle-high", hertz_dependency_name, hertz_module_name, "v0.9.1", "", "1.18", "", []string{"go", "build", "test_hertz_basic.go", "basic_func.go"}))
}

func TestBasicHertz(t *testing.T, env ...string) {
	UseApp("hertz/v0.9.0")
	RunGoBuild(t, "go", "build", "test_hertz_basic.go", "basic_func.go")
	RunApp(t, "test_hertz_basic", env...)
}

func TestBasicHertzWithHook(t *testing.T, env ...string) {
	UseApp("hertz/v0.9.0")
	RunGoBuild(t, "go", "build", "test_hertz_with_hook.go", "basic_func.go")
	RunApp(t, "test_hertz_with_hook", env...)
}

func TestBasicHertzWithException(t *testing.T, env ...string) {
	UseApp("hertz/v0.9.0")
	RunGoBuild(t, "go", "build", "test_hertz_with_exception.go", "basic_func.go")
	RunApp(t, "test_hertz_with_exception", env...)
}

func TestBasicHertzWithRegex(t *testing.T, env ...string) {
	UseApp("hertz/v0.9.0")
	RunGoBuild(t, "go", "build", "test_hertz_with_regex.go", "basic_func.go")
	RunApp(t, "test_hertz_with_regex", env...)
}
