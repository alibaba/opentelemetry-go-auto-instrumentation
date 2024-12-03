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
		NewGeneralTestCase("logrus-test", "logrus", "", "", "1.21", "", TestLogrus),
	)
}

func TestLogrus(t *testing.T, env ...string) {
	UseApp("logrus")
	RunGoBuild(t, "go", "build", "test_logrus.go", "http_server.go")
	_, stderr := RunApp(t, "test_logrus", env...)
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
