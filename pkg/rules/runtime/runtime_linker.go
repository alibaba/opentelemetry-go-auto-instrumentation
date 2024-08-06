//go:build ignore

package runtime

import (
	_ "unsafe"
)

//go:linkname ot_get_trace_context_from_gls ot_get_trace_context_from_gls
var ot_get_trace_context_from_gls = _ot_gls_get_trace_context_impl

//go:linkname ot_get_baggage_container_from_gls ot_get_baggage_container_from_gls
var ot_get_baggage_container_from_gls = _ot_gls_get_baggage_container_impl

//go:linkname ot_set_trace_context_to_gls ot_set_trace_context_to_gls
var ot_set_trace_context_to_gls = _ot_gls_set_trace_context_impl

//go:linkname ot_set_baggage_container_to_gls ot_set_baggage_container_to_gls
var ot_set_baggage_container_to_gls = _ot_gls_set_baggage_container_impl

//go:nosplit
func _ot_gls_get_trace_context_impl() interface{} {
	return getg().m.curg.ot_trace_context
}

//go:nosplit
func _ot_gls_get_baggage_container_impl() interface{} {
	return getg().m.curg.ot_baggage_container
}

//go:nosplit
func _ot_gls_set_trace_context_impl(v interface{}) {
	getg().m.curg.ot_trace_context = v
}

//go:nosplit
func _ot_gls_set_baggage_container_impl(v interface{}) {
	getg().m.curg.ot_baggage_container = v
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
