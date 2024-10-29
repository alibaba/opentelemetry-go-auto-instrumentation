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

package runtime

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	// Add tls field to struct "g" in runtime2.go
	api.NewStructRule("runtime", "g", "otel_trace_context", "interface{}").
		Register()
	api.NewStructRule("runtime", "g", "otel_baggage_container", "interface{}").
		Register()

	// This solely inspired by skywalking-go
	// https://github.com/apache/skywalking-go/blob/5d7bd5e8e435ec5ab1a61793cd08e6a403893a55/tools/go-agent/instrument/runtime/instrument.go#L75
	api.NewRule("runtime",
		"newproc1", "", "defer func(){ retVal0.otel_trace_context = contextPropagate(callergp.otel_trace_context); retVal0.otel_baggage_container = contextPropagate(callergp.otel_baggage_container); }()", "").
		WithUseRaw(true).
		Register()

	api.NewFileRule("runtime", "runtime_linker.go").
		Register()
}
