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
