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
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const HelloworldAppName = "helloworld"

func TestRunHelloworld(t *testing.T) {
	UseApp(HelloworldAppName)

	RunSet(t, "-rule=")
	RunGoBuild(t, "go", "build") // no test rules, build as usual
	stdout, _ := RunApp(t, HelloworldAppName)
	ExpectContains(t, stdout, "helloworld")

	RunSet(t, UseTestRules("test_fmt.json"))
	RunGoBuild(t, "go", "build")
	stdout, stderr := RunApp(t, HelloworldAppName)
	ExpectContains(t, stdout, "olleH")
	ExpectContains(t, stderr, "Entering hook1") // println writes to stderr
	ExpectContains(t, stderr, "Exiting hook1")
	ExpectContains(t, stderr, "555")
	ExpectContains(t, stderr, "internalFn")
	ExpectContains(t, stderr, "7632")
	ExpectContains(t, stderr, "init")
	ExpectContains(t, stderr, "init2")
	ExpectContains(t, stderr, "30258") //0x7632
	ExpectContains(t, stderr, "GOOD")
	ExpectNotContains(t, stderr, "BAD")
	ExpectContains(t, stderr, "GCMG")
	ExpectContains(t, stderr, "BYD")

	text := ReadInstrumentLog(t, filepath.Join("fmt", "print.go"))
	re := regexp.MustCompile(".*OtelOnEnterTrampoline.*OtelOnExitTrampoline.*")
	matches := re.FindAllString(text, -1)
	if len(matches) < 1 {
		t.Fatalf("expecting at least one match")
	}
}

func runModVendor(t *testing.T) {
	cmd := exec.Command("go", "mod", "tidy")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go mod tidy failed: %v\n%s", err, out)
	}
	cmd = exec.Command("go", "mod", "vendor")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go mod vendor failed: %v\n%s", err, out)
	}
}

func TestBuildHelloworldWithVendor1(t *testing.T) {
	UseApp(HelloworldAppName)
	runModVendor(t)
	RunGoBuild(t, "go", "build")
	ExpectPreprocessNotContains(t, util.DebugLogFile, "Bad match")
}

func TestBuildHelloworldWithVendor2(t *testing.T) {
	UseApp(HelloworldAppName)
	runModVendor(t)
	RunGoBuild(t, "go", "build", "-mod=vendor")
	ExpectPreprocessNotContains(t, util.DebugLogFile, "Bad match")
}

func TestBuildHelloworldWithVendor3(t *testing.T) {
	UseApp(HelloworldAppName)
	runModVendor(t)
	RunGoBuild(t, "go", "build", "-mod", "vendor")
	ExpectPreprocessNotContains(t, util.DebugLogFile, "Bad match")
}

func TestBuildHelloworldWithVendor4(t *testing.T) {
	UseApp(HelloworldAppName)
	runModVendor(t)
	RunGoBuild(t, "go", "build", "-mod=mod")
	ExpectPreprocessNotContains(t, util.DebugLogFile, "run go mod vendor")
	ExpectPreprocessNotContains(t, util.DebugLogFile, "Bad match")
}
