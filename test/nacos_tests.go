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
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"testing"
	"time"
)

const nacos_dependency_name = "github.com/gomodule/redigo"
const nacos_module_name = "nacos"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("nacos-2.0.0-config-test", nacos_module_name, "v2.0.0", "v2.1.0", "1.18", "", TestNacos200Config),
		NewGeneralTestCase("nacos-2.0.0-service-test", nacos_module_name, "v2.0.0", "v2.1.0", "1.18", "", TestNacos200Service),
		NewGeneralTestCase("nacos-2.1.0-config-test", nacos_module_name, "v2.1.0", "", "1.18", "", TestNacos210Config),
		NewGeneralTestCase("nacos-2.1.0-service-test", nacos_module_name, "v2.1.0", "", "1.18", "", TestNacos210Service))
}

func TestNacos200Config(t *testing.T, env ...string) {
	nacosC, nacosPort := initNacosContainer()
	defer clearNacosContainer(nacosC)
	UseApp("nacos/v2.0.0")
	RunGoBuild(t, "go", "build", "test_nacos_config.go")
	env = append(env, "NACOS_PORT="+nacosPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_NACOS_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_nacos_config", env...)
}

func TestNacos200Service(t *testing.T, env ...string) {
	nacosC, nacosPort := initNacosContainer()
	defer clearNacosContainer(nacosC)
	UseApp("nacos/v2.0.0")
	RunGoBuild(t, "go", "build", "test_nacos_service.go")
	env = append(env, "NACOS_PORT="+nacosPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_NACOS_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_nacos_service", env...)
}

func TestNacos210Config(t *testing.T, env ...string) {
	nacosC, nacosPort := initNacosContainer()
	defer clearNacosContainer(nacosC)
	UseApp("nacos/v2.1.0")
	RunGoBuild(t, "go", "build", "test_nacos_config.go")
	env = append(env, "NACOS_PORT="+nacosPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_NACOS_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_nacos_config", env...)
}

func TestNacos210Service(t *testing.T, env ...string) {
	// start nacos image
	nacosC, nacosPort := initNacosContainer()
	defer clearNacosContainer(nacosC)
	UseApp("nacos/v2.1.0")
	RunGoBuild(t, "go", "build", "test_nacos_service.go")
	env = append(env, "NACOS_PORT="+nacosPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_NACOS_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_nacos_service", env...)
}

func initNacosContainer() (testcontainers.Container, nat.Port) {
	req := testcontainers.ContainerRequest{
		Image:        "nacos/nacos-server:latest",
		ExposedPorts: []string{"8848/tcp", "9848/tcp"},
		Env: map[string]string{
			"MODE": "standalone",
		},
	}
	nacosC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	port, err := nacosC.MappedPort(context.Background(), "8848")
	if err != nil {
		panic(err)
	}
	return nacosC, port
}

func clearNacosContainer(nacosC testcontainers.Container) {
	if err := nacosC.Terminate(context.Background()); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
}
