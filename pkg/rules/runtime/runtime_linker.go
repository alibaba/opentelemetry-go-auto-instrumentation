// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
