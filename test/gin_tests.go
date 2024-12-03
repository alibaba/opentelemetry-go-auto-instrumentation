// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("gin-test-html", "gin", "", "", "1.21", "", TestGinHTML),
		NewGeneralTestCase("gin-test-pattern", "gin", "", "", "1.21", "", TestGinPattern),
	)
}

func TestGinHTML(t *testing.T, env ...string) {
	UseApp("gin")
	RunGoBuild(t, "go", "build", "test_gin_html.go")
	RunApp(t, "test_gin_html", env...)
}

func TestGinPattern(t *testing.T, env ...string) {
	UseApp("gin")
	RunGoBuild(t, "go", "build", "test_gin_pattern.go")
	RunApp(t, "test_gin_pattern", env...)
}
