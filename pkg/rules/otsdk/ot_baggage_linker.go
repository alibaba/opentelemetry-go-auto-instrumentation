//go:build ignore

package baggage

import (
	_ "unsafe"
)

var (
	GetBaggageContainerFromGLS = func() interface{} { return nil }
	SetBaggageContainerToGLS   = func(interface{}) {}
)

//go:linkname ot_get_baggage_container_from_gls ot_get_baggage_container_from_gls
var ot_get_baggage_container_from_gls func() interface{}

//go:linkname ot_set_baggage_container_to_gls ot_set_baggage_container_to_gls
var ot_set_baggage_container_to_gls func(interface{})

func init() {
	if ot_get_baggage_container_from_gls != nil && ot_set_baggage_container_to_gls != nil {
		GetBaggageContainerFromGLS = ot_get_baggage_container_from_gls
		SetBaggageContainerToGLS = ot_set_baggage_container_to_gls
	}
}
