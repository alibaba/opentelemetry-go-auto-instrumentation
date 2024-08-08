package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("nethttp-basic-test", "", "databasesql", "", "", "1.18", "", TestBasicNetHttp),
	)
}

func TestBasicNetHttp(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunInstrument(t, "-debuglog", "--", "test_http.go", "http_server.go")
	RunApp(t, "test_http", env...)
}
