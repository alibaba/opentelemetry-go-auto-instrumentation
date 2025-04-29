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
	"time"

	_ "go.opentelemetry.io/otel"
	"golang.org/x/time/rate"

	// When building with vendor mode, we requires the following packages
	// to be imported to ensure they are included in the vendor directory.
	// In this way, we can build the artifact offline, i.e. without the
	// need to download the dependencies.
	_ "github.com/mohae/deepcopy"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	_ "go.opentelemetry.io/contrib/instrumentation/runtime"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	_ "go.opentelemetry.io/otel/exporters/prometheus"
	_ "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	_ "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	_ "go.opentelemetry.io/otel/exporters/zipkin"
	_ "go.opentelemetry.io/otel/sdk/trace/tracetest"
	_ "go.opentelemetry.io/otel/semconv/v1.19.0"
	_ "go.opentelemetry.io/otel/semconv/v1.30.0"
	_ "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/runtime/protoimpl"
)

func main() {
	n, _ := fmt.Printf("helloworld%s", "ingodwetrust")
	println(n)

	println(rate.Every(time.Duration(1) * time.Second))
}
