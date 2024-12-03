// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	_ "unsafe"
)

var (
	GetTraceContextFromGLS   = func() interface{} { return nil }
	SetTraceContextToGLS     = func(interface{}) {}
	SetBaggageContainerToGLS = func(interface{}) {}
)

//go:linkname otel_get_trace_context_from_gls otel_get_trace_context_from_gls
var otel_get_trace_context_from_gls func() interface{}

//go:linkname otel_set_trace_context_to_gls otel_set_trace_context_to_gls
var otel_set_trace_context_to_gls func(interface{})

//go:linkname otel_set_baggage_container_to_gls otel_set_baggage_container_to_gls
var otel_set_baggage_container_to_gls func(interface{})

func init() {
	if otel_get_trace_context_from_gls != nil && otel_set_trace_context_to_gls != nil {
		GetTraceContextFromGLS = otel_get_trace_context_from_gls
		SetTraceContextToGLS = otel_set_trace_context_to_gls
	}
	if otel_set_baggage_container_to_gls != nil {
		SetBaggageContainerToGLS = otel_set_baggage_container_to_gls
	}
}
