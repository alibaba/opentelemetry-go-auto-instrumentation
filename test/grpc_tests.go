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

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("grpc-basic-test", "grpc", "", "", "1.21", "", TestBasicGrpc),
		NewGeneralTestCase("grpc-stream-test", "grpc", "", "", "1.21", "", TestGrpcStream),
	)
}

func TestBasicGrpc(t *testing.T, env ...string) {
	UseApp("grpc")
	RunInstrument(t, "-debuglog", "--", "test_grpc_basic.go", "server.go", "grpc_server.go", "grpc.pb.go", "grpc_grpc.pb.go")
	RunApp(t, "test_grpc_basic", env...)
}

func TestGrpcStream(t *testing.T, env ...string) {
	UseApp("grpc_stream")
	RunInstrument(t, "-debuglog", "--", "test_grpc_stream.go")
	RunApp(t, "test_grpc_stream", env...)
}
