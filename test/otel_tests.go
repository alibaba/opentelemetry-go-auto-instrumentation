// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("otel-span-from-context-test", "otel", "", "", "1.18", "", TestSpanFromContext),
	)
}

func TestSpanFromContext(t *testing.T, env ...string) {
	UseApp("otel")
	RunGoBuild(t, "go", "build", "--", "test_span_from_context.go")
	stdout, _ := RunApp(t, "test_span_from_context", env...)
	ExpectContains(t, stdout, "GET /otel")
}
