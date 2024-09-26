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
