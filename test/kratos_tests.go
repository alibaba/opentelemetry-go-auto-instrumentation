// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

const kratos_dependency_name = "github.com/go-kratos/kratos/v2"
const kratos_module_name = "kratos"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("kratos-http-test", kratos_module_name, "", "", "1.18", "", TestKratosHttp),
		NewGeneralTestCase("kratos-grpc-test", kratos_module_name, "", "", "1.18", "", TestKratosGrpc),
		NewLatestDepthTestCase("kratos-latest-depth-grpc", kratos_dependency_name, kratos_module_name, "v2.6.3", "", "1.22", "", TestKratosGrpc),
		NewLatestDepthTestCase("kratos-latest-depth-http", kratos_dependency_name, kratos_module_name, "v2.6.3", "", "1.22", "", TestKratosHttp),
		NewMuzzleTestCase("kratos-muzzle-grpc", kratos_dependency_name, kratos_module_name, "v2.6.3", "", "1.22", "", []string{"go", "build", "test_kratos_grpc.go", "server.go"}),
		NewMuzzleTestCase("kratos-muzzle-http", kratos_dependency_name, kratos_module_name, "v2.6.3", "", "1.22", "", []string{"go", "build", "test_kratos_http.go", "server.go"}),
	)
}

func TestKratosGrpc(t *testing.T, env ...string) {
	UseApp("kratos/v2.6.3")
	RunGoBuild(t, "go", "build", "test_kratos_grpc.go", "server.go")
	env = append(env, "OTEL_INSTRUMENTATION_KRATOS_EXPERIMENTAL_SPAN_ATTRIBUTES=true")
	RunApp(t, "test_kratos_grpc", env...)
}

func TestKratosHttp(t *testing.T, env ...string) {
	UseApp("kratos/v2.6.3")
	RunGoBuild(t, "go", "build", "test_kratos_http.go", "server.go")
	env = append(env, "OTEL_INSTRUMENTATION_KRATOS_EXPERIMENTAL_SPAN_ATTRIBUTES=true")
	RunApp(t, "test_kratos_http", env...)
}
