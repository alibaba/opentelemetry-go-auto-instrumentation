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

package main

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"net/http"
	"strings"
	"time"
)

func main() {
	go func() {
		startup()
	}()
	time.Sleep(5 * time.Second)
	client := http.Client{}
	client.Get("http://localhost:8000/helloworld/kratos")
	fmt.Printf("Send http request to kratos")

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		for i, _ := range stubs[0] {
			span := stubs[0][i]
			println(span.Name)
			for _, attr := range span.Attributes {
				fmt.Printf("%v %v\n", attr.Key, attr.Value)
			}
			println()
		}

		verifier.NewSpanVerifier().
			HasStringAttribute("kratos.protocol.type", "http").
			HasStringAttribute("kratos.service.name", "opentelemetry-kratos-server").
			HasStringAttribute("kratos.service.id", "opentelemetry-id").
			HasStringAttribute("kratos.service.version", "v1").
			HasStringAttribute("kratos.service.meta.agent", "opentelemetry-go").
			HasItemInStringSliceAttribute("kratos.service.endpoint", 0, func(s string) (bool, string) {
				return strings.Contains(s, ":9000"), fmt.Sprintf("First endpoint should be xxx:9000, actual value: %v", s)
			}).
			HasItemInStringSliceAttribute("kratos.service.endpoint", 1, func(s string) (bool, string) {
				return strings.Contains(s, ":8000"), fmt.Sprintf("First endpoint should be xxx:8000, actual value: %v", s)
			}).
			Verify(stubs[0][2])
	}, 1)
}
