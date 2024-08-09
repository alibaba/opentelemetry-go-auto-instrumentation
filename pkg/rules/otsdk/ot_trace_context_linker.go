//go:build ignore

package trace

import (
	_ "unsafe"
)

var (
	GetTraceContextFromGLS   = func() interface{} { return nil }
	SetTraceContextToGLS     = func(interface{}) {}
	SetBaggageContainerToGLS = func(interface{}) {}
)

//go:linkname ot_get_trace_context_from_gls ot_get_trace_context_from_gls
var ot_get_trace_context_from_gls func() interface{}

//go:linkname ot_set_trace_context_to_gls ot_set_trace_context_to_gls
var ot_set_trace_context_to_gls func(interface{})

//go:linkname ot_set_baggage_container_to_gls ot_set_baggage_container_to_gls
var ot_set_baggage_container_to_gls func(interface{})

func init() {
	if ot_get_trace_context_from_gls != nil && ot_set_trace_context_to_gls != nil {
		GetTraceContextFromGLS = ot_get_trace_context_from_gls
		SetTraceContextToGLS = ot_set_trace_context_to_gls
	}
	if ot_set_baggage_container_to_gls != nil {
		SetBaggageContainerToGLS = ot_set_baggage_container_to_gls
	}
}
