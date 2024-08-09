package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("nethttp-basic-test", "", "nethttp", "", "", "1.18", "", TestBasicNetHttp),
		NewGeneralTestCase("nethttp-http-2-test", "", "nethttp", "", "", "1.18", "", TestHttp2),
		NewGeneralTestCase("nethttp-https-test", "", "nethttp", "", "", "1.18", "", TestHttps),
	)
}

func TestBasicNetHttp(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunInstrument(t, "-debuglog", "--", "test_http.go", "http_server.go")
	RunApp(t, "test_http", env...)
}

func TestHttp2(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunInstrument(t, "-debuglog", "--", "test_http_2.go", "http_server.go")
	RunApp(t, "test_http_2", env...)
}

func TestHttps(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunInstrument(t, "-debuglog", "--", "test_https.go", "http_server.go")
	RunApp(t, "test_https", env...)
}
