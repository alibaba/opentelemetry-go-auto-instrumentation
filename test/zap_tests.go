// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bufio"
	"strings"
	"testing"
)

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("zap-test", "zap", "", "", "1.21", "", TestZap),
	)
}

func TestZap(t *testing.T, env ...string) {
	UseApp("zap")
	RunGoBuild(t, "go", "build", "test_zap.go", "http_server.go")
	_, stderr := RunApp(t, "test_zap", env...)
	reader := strings.NewReader(stderr)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[test debugging]") {
			continue
		}
		ExpectContains(t, line, "trace_id")
		ExpectContains(t, line, "span_id")
	}
}
