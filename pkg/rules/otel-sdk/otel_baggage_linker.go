// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage

import (
	_ "unsafe"
)

var (
	GetBaggageContainerFromGLS = func() interface{} { return nil }
	SetBaggageContainerToGLS   = func(interface{}) {}
)

//go:linkname otel_get_baggage_container_from_gls otel_get_baggage_container_from_gls
var otel_get_baggage_container_from_gls func() interface{}

//go:linkname otel_set_baggage_container_to_gls otel_set_baggage_container_to_gls
var otel_set_baggage_container_to_gls func(interface{})

func init() {
	if otel_get_baggage_container_from_gls != nil && otel_set_baggage_container_to_gls != nil {
		GetBaggageContainerFromGLS = otel_get_baggage_container_from_gls
		SetBaggageContainerToGLS = otel_set_baggage_container_to_gls
	}
}
