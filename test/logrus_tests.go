// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	RunInstrument(t, "-debuglog", "go", "build", "test_logrus.go", "http_server.go")
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
