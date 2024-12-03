// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

const fasthttp_dependency_name = "github.com/valyala/fasthttp"
const fasthttp_module_name = "fasthttp"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("basic-fasthttp-test", fasthttp_module_name, "", "", "1.18", "", TestBasicFastHttp),
		NewGeneralTestCase("basic-fasthttps-test", fasthttp_module_name, "", "", "1.18", "", TestBasicFastHttps),
		NewLatestDepthTestCase("fasthttp-latestdepth", fasthttp_dependency_name, fasthttp_module_name, "v1.45.0", "", "1.18", "", TestBasicFastHttp),
		NewMuzzleTestCase("fasthttp-muzzle", fasthttp_dependency_name, fasthttp_module_name, "v1.45.0", "", "1.18", "", []string{"go", "build", "test_basic_http.go", "server.go"}))
}

func TestBasicFastHttp(t *testing.T, env ...string) {
	UseApp("fasthttp/v1.45.0")
	RunGoBuild(t, "go", "build", "test_basic_http.go", "server.go")
	RunApp(t, "test_basic_http", env...)
}

func TestBasicFastHttps(t *testing.T, env ...string) {
	UseApp("fasthttp/v1.45.0")
	RunGoBuild(t, "go", "build", "test_basic_https.go", "server.go")
	RunApp(t, "test_basic_https", env...)
}
