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
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"time"
)

const es_v8_dependency_name = "github.com/elastic/go-elasticsearch/v8"
const es_v8_module_name = "elasticsearch"

const defaultHTTPPort = "9200"
const defaultTCPPort = "9300"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("es-crud-test", es_v8_module_name, "v8.0.0", "v8.15.1", "1.18", "", TestESCrud),
		NewLatestDepthTestCase("es-crud-latestdepth-test", es_v8_dependency_name, es_v8_module_name, "v8.0.0", "v8.15.1", "1.18", "", TestESCrud),
		NewMuzzleTestCase("es-muzzle", es_v8_dependency_name, es_v8_module_name, "v8.0.0", "v8.15.1", "1.18", "", []string{"go", "build", "test_es_crud.go"}),
	)
}

func TestESCrud(t *testing.T, env ...string) {
	esC, esPort := initElasticSearchContainer()
	defer testcontainers.CleanupContainer(t, esC)
	UseApp("elasticsearch/v8.0.0")
	RunInstrument(t, "-debuglog", "go", "build", "test_es_crud.go")
	env = append(env, "OTEL_ES_PORT="+esPort.Port())
	RunApp(t, "test_es_crud", env...)
}

func initElasticSearchContainer() (testcontainers.Container, nat.Port) {
	ctx := context.Background()
	elasticsearchContainer, err := runElasticSearchContainer(ctx)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	port, err := elasticsearchContainer.MappedPort(context.Background(), "9200")
	if err != nil {
		panic(err)
	}
	return elasticsearchContainer, port
}

func runElasticSearchContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.elastic.co/elasticsearch/elasticsearch:8.9.0",
			Env: map[string]string{
				"discovery.type": "single-node",
				"cluster.routing.allocation.disk.threshold_enabled": "false",
			},
			ExposedPorts: []string{
				defaultHTTPPort + "/tcp",
				defaultTCPPort + "/tcp",
			},
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, fmt.Errorf("generic container: %w", err)
	}

	return container, nil
}
