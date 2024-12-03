// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

const grpc_dependency_name = "google.golang.org/grpc"
const grpc_module_name = "grpc"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("grpc-basic-test", grpc_module_name, "v1.44.0", "", "1.21", "", TestBasicGrpc),
		NewGeneralTestCase("grpc-fail-status-test", grpc_module_name, "v1.44.0", "", "1.21", "", TestGrpcStatus),
		NewGeneralTestCase("grpc-stream-test", grpc_module_name, "v1.44.0", "", "1.21", "", TestGrpcStream),
		NewLatestDepthTestCase("grpc-latest-depth", grpc_dependency_name, grpc_module_name, "v1.44.0", "", "1.21", "", TestBasicGrpc),
		NewMuzzleTestCase("grpc-muzzle", grpc_dependency_name, grpc_module_name, "v1.44.0", "", "1.21", "", []string{"go", "build", "test_grpc_basic.go", "grpc_common.go", "grpc.pb.go", "grpc_grpc.pb.go"}),
	)
}

func TestBasicGrpc(t *testing.T, env ...string) {
	UseApp("grpc/v1.44.0")
	RunGoBuild(t, "go", "build", "test_grpc_basic.go", "grpc_common.go", "grpc.pb.go", "grpc_grpc.pb.go")
	env = append(env, "GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn")
	RunApp(t, "test_grpc_basic", env...)
}

func TestGrpcStatus(t *testing.T, env ...string) {
	UseApp("grpc/v1.44.0")
	RunGoBuild(t, "go", "build", "test_grpc_fail_status.go", "grpc_common.go", "grpc.pb.go", "grpc_grpc.pb.go")
	env = append(env, "GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn")
	RunApp(t, "test_grpc_fail_status", env...)
}

func TestGrpcStream(t *testing.T, env ...string) {
	UseApp("grpc_stream")
	RunGoBuild(t, "go", "build", "test_grpc_stream.go")
	env = append(env, "GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn")
	RunApp(t, "test_grpc_stream", env...)
}
