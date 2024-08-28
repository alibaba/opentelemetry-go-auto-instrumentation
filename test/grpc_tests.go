package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("grpc-test", "grpc", "", "", "1.21", "", TestBasicGrpc),
	)
}

func TestBasicGrpc(t *testing.T, env ...string) {
	UseApp("grpc")
	RunInstrument(t, "-debuglog", "--", "grpc.go", "server.go", "grpc_server.go", "grpc.pb.go", "grpc_grpc.pb.go")
	RunApp(t, "grpc", env...)
}
