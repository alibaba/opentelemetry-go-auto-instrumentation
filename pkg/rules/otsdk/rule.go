// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otsdk

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "ot_trace_context_linker.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "trace-context/ot_trace_context.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "trace-context/span.go").WithReplace(true).Register()
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "trace-context/tracer.go").WithReplace(true).Register()
	api.NewFileRule("go.opentelemetry.io/otel", "trace-context/trace.go").WithReplace(true).Register()
	// baggage
	api.NewFileRule("go.opentelemetry.io/otel/baggage", "ot_baggage_linker.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/baggage", "ot_baggage_util.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/baggage", "baggage/context.go").WithReplace(true).Register()
}
