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
	"github.com/testcontainers/testcontainers-go/modules/elasticsearch"
	"testing"
	"time"
)

const es_v8_dependency_name = "github.com/elastic/go-elasticsearch/v8"
const es_v8_module_name = "elasticsearch"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("es-crud-test", es_v8_module_name, "v8.0.0", "v8.15.1", "1.18", "", TestESCrud),
		NewLatestDepthTestCase("es-crud-latestdepth-test", es_v8_dependency_name, es_v8_module_name, "v8.0.0", "v8.15.1", "1.18", "", TestESCrud),
		NewMuzzleTestCase("es-muzzle", es_v8_dependency_name, es_v8_module_name, "v8.0.0", "v8.15.1", "1.18", "", []string{"test_es_crud.go"}),
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
	elasticsearchContainer, err := elasticsearch.Run(ctx, "docker.elastic.co/elasticsearch/elasticsearch:8.9.0")
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	port, err := elasticsearchContainer.MappedPort(context.Background(), "6379")
	if err != nil {
		panic(err)
	}
	return elasticsearchContainer, port
}
