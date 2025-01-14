// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

const gomicro_dependency_name = "go-micro.dev/v5"
const gomicro_module_name = "gomicro"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("gomicro-basic-test", gomicro_module_name, "v5.0.0", "", "1.21", "", TestBasicGoMicro),
	)
}

func TestBasicGoMicro(t *testing.T, env ...string) {
	UseApp("gomicro/v5.3.0")
	RunGoBuild(t, "go", "build", "test_gomicro.go", "test_goserver.go")
	env = append(env, "GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn")
	RunApp(t, "test_gomicro", env...)
}
