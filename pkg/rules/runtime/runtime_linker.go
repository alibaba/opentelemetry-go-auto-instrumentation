// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	_ "unsafe"
)

//go:linkname otel_get_trace_context_from_gls otel_get_trace_context_from_gls
var otel_get_trace_context_from_gls = _otel_gls_get_trace_context_impl

//go:linkname otel_get_baggage_container_from_gls otel_get_baggage_container_from_gls
var otel_get_baggage_container_from_gls = _otel_gls_get_baggage_container_impl

//go:linkname otel_set_trace_context_to_gls otel_set_trace_context_to_gls
var otel_set_trace_context_to_gls = _otel_gls_set_trace_context_impl

//go:linkname otel_set_baggage_container_to_gls otel_set_baggage_container_to_gls
var otel_set_baggage_container_to_gls = _otel_gls_set_baggage_container_impl

//go:nosplit
func _otel_gls_get_trace_context_impl() interface{} {
	return getg().m.curg.otel_trace_context
}

//go:nosplit
func _otel_gls_get_baggage_container_impl() interface{} {
	return getg().m.curg.otel_baggage_container
}

//go:nosplit
func _otel_gls_set_trace_context_impl(v interface{}) {
	getg().m.curg.otel_trace_context = v
}

//go:nosplit
func _otel_gls_set_baggage_container_impl(v interface{}) {
	getg().m.curg.otel_baggage_container = v
}

type ContextSnapshoter interface {
	TakeSnapShot() interface{}
}

func contextPropagate(tls interface{}) interface{} {
	if tls == nil {
		return nil
	}
	if taker, ok := tls.(ContextSnapshoter); ok {
		return taker.TakeSnapShot()
	}
	return tls
}
