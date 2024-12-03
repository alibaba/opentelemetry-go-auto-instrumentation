// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("nethttp-basic-test", "nethttp", "", "", "1.18", "", TestBasicNetHttp),
		NewGeneralTestCase("nethttp-http-2-test", "nethttp", "", "", "1.18", "", TestHttp2),
		NewGeneralTestCase("nethttp-https-test", "nethttp", "", "", "1.18", "", TestHttps),
		NewGeneralTestCase("nethttp-metric-test", "nethttp", "", "", "1.18", "", TestHttpMetric),
	)
}

func TestBasicNetHttp(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunGoBuild(t, "go", "build", "test_http.go", "http_server.go")
	RunApp(t, "test_http", env...)
}

func TestHttp2(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunGoBuild(t, "go", "build", "test_http_2.go", "http_server.go")
	RunApp(t, "test_http_2", env...)
}

func TestHttps(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunGoBuild(t, "go", "build", "test_https.go", "http_server.go")
	RunApp(t, "test_https", env...)
}

func TestHttpMetric(t *testing.T, env ...string) {
	UseApp("nethttp")
	RunGoBuild(t, "go", "build", "test_http_metrics.go", "http_server.go")
	RunApp(t, "test_http_metrics", env...)
}
