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

package otel_context

import (
	_ "unsafe"

	"github.com/alibaba/loongsuite-go/pkg/api"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	spectrace "go.opentelemetry.io/otel/trace"
)

//go:linkname spanFromContextOnExit go.opentelemetry.io/otel/trace.spanFromContextOnExit
func spanFromContextOnExit(call api.CallContext, span spectrace.Span) {
	if !span.SpanContext().IsValid() {
		// Only replace with a GLS span when one exists. SpanFromGLS() returns
		// nil when the goroutine TraceContext is unset or empty; substituting
		// that nil would break callers such as SpanContextFromContext that
		// assume SpanFromContext never returns a nil Span (OTel returns a
		// non-nil noop span instead).
		if glsSpan := sdktrace.SpanFromGLS(); glsSpan != nil {
			call.SetReturnVal(0, glsSpan)
		}
	}
	return
}
