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
	"regexp"
	"testing"
)

const HelloworldAppName = "helloworld"

func TestRunHelloworld(t *testing.T) {
	UseApp(HelloworldAppName)

	RunInstrument(t, "-debuglog")
	stdout, _ := RunApp(t, HelloworldAppName)
	ExpectContains(t, stdout, "helloworld")

	RunInstrument(t, "-debuglog", "-disablerules=") // use fmt rules
	RunInstrument(t, "-restore")                    // restore then
	RunApp(t, HelloworldAppName)                    // run the app again
	ExpectContains(t, stdout, "helloworld")         // nothing should change

	RunInstrument(t, "-debuglog", "-disablerules=")
	stdout, stderr := RunApp(t, HelloworldAppName)
	ExpectContains(t, stdout, "olleH")
	ExpectContains(t, stderr, "Entering hook1") // println writes to stderr
	ExpectContains(t, stderr, "Exiting hook1")
	ExpectContains(t, stderr, "555")
	ExpectContains(t, stderr, "internalFn")
	ExpectContains(t, stderr, "GCMG")
	ExpectContains(t, stderr, "7632")
	ExpectContains(t, stderr, "init")
	ExpectContains(t, stderr, "init2")
	ExpectContains(t, stderr, "512")
	ExpectContains(t, stderr, "30258") //0x7632
	ExpectContains(t, stderr, "GOOD")
	ExpectNotContains(t, stderr, "BAD")

	text := ReadInstrumentLog(t, "debug_fn_print.go")
	re := regexp.MustCompile(".*OtelOnEnterTrampoline.*OtelOnExitTrampoline.*")
	matches := re.FindAllString(text, -1)
	if len(matches) < 1 {
		t.Fatalf("expecting at least one match")
	}
}

func TestBuildHelloworldWithVendor1(t *testing.T) {
	UseApp(HelloworldAppName)
	RunCmd([]string{"go", "mod", "vendor"})
	RunInstrument(t, "-debuglog")
}

func TestBuildHelloworldWithVendor2(t *testing.T) {
	UseApp(HelloworldAppName)
	RunCmd([]string{"go", "mod", "vendor"})
	RunInstrument(t, "-debuglog", "--", "-mod=vendor")
}

func TestBuildHelloworldWithVendor3(t *testing.T) {
	UseApp(HelloworldAppName)
	RunCmd([]string{"go", "mod", "vendor"})
	RunInstrument(t, "-debuglog", "--", "-mod", "vendor")
}
