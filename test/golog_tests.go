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
		NewGeneralTestCase("golog-test", "golog", "", "", "1.18", "", TestGoLog),
	)
}

func TestGoLog(t *testing.T, env ...string) {
	UseApp("golog")
	RunGoBuild(t, "go", "build", "test_glog.go")
	_, stderr := RunApp(t, "test_glog", env...)
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
