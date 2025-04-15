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

package core

// @@ There is a scenario where we want users to rely on our package (alibaba/pkg/core),
// and after executing go mod vendor, all related dependencies are fetched
// locally. The subsequent build process can then be performed offline,
// without relying on go mod tidy/vendor to fetch any additional dependencies.
// Therefore, we need to proactively list all hook-related code and their
// potential dependencies in advance and import them here.
// This makes no sense for most users, but it does make sense for users whose
// build machine is not connected to the internet.
import (
	_ "github.com/mohae/deepcopy"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	_ "go.opentelemetry.io/contrib/instrumentation/runtime"
	_ "go.opentelemetry.io/otel"
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
	_ "go.opentelemetry.io/otel/semconv/v1.30.0"
	_ "google.golang.org/protobuf/proto"
)
