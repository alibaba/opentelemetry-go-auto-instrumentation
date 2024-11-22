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

import "testing"

const fiberV2_dependency_name = "github.com/gofiber/fiber/v2"
const fiberV2_module_name = "fiberV2"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("basic-fiberV2-test", fiberV2_module_name, "", "", "1.18", "", TestBasicFiberV2),
		NewGeneralTestCase("basic-fiberV2s-test", fiberV2_module_name, "", "", "1.18", "", TestBasicFiberV2Https),
		NewLatestDepthTestCase("fiberV2-latestdepth", fiberV2_dependency_name, fiberV2_module_name, "v2.43.0", "", "1.18", "", TestBasicFiberV2),
		NewMuzzleTestCase("fiberV2-muzzle", fiberV2_dependency_name, fiberV2_module_name, "v2.43.0", "", "1.18", "", []string{"go", "build", "fiber_http.go"}))
}

func TestBasicFiberV2(t *testing.T, env ...string) {
	UseApp("fiberV2/v2.52.5")
	RunInstrument(t, "-debuglog", "go", "build", "fiber_http.go")
	RunApp(t, "fiber_http", env...)
}

func TestBasicFiberV2Https(t *testing.T, env ...string) {
	UseApp("fiberV2/v2.52.5")
	RunInstrument(t, "-debuglog", "go", "build", "fiber_https.go")
	RunApp(t, "fiber_https", env...)
}
