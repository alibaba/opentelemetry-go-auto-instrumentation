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

package trace

import (
	"fmt"

	trace "go.opentelemetry.io/otel/trace"
)

const maxSpans = 300

type traceContext struct {
	sw    *spanWrapper
	n     int
	DataK []string
	DataV []interface{}
	lcs   trace.Span
}

type spanWrapper struct {
	span trace.Span
	prev *spanWrapper
}

func (tc *traceContext) size() int {
	return tc.n
}

func (tc *traceContext) add(span trace.Span) bool {
	if tc.n > 0 {
		if tc.n >= maxSpans {
			return false
		}
	}
	wrapper := &spanWrapper{span, tc.sw}
	// local root span
	if tc.n == 0 {
		tc.lcs = span
	}
	tc.sw = wrapper
	tc.n++
	return true
}

//go:norace
func (tc *traceContext) tail() trace.Span {
	if tc.n == 0 {
		return nil
	} else {
		return tc.sw.span
	}
}

func (tc *traceContext) localRootSpan() trace.Span {
	if tc.n == 0 {
		return nil
	} else {
		return tc.lcs
	}
}

func (tc *traceContext) del(span trace.Span) {
	if tc.n == 0 {
		return
	}
	addr := &tc.sw
	cur := tc.sw
	for cur != nil {
		sc1 := cur.span.SpanContext()
		sc2 := span.SpanContext()
		if sc1.TraceID() == sc2.TraceID() && sc1.SpanID() == sc2.SpanID() {
			*addr = cur.prev
			tc.n--
			break
		}
		addr = &cur.prev
		cur = cur.prev
	}
}

func (tc *traceContext) clear() {
	tc.sw = nil
	tc.n = 0
	tc.DataK = make([]string, 0)
	tc.DataV = make([]interface{}, 0)
	SetBaggageContainerToGLS(nil)
}

//go:norace
func (tc *traceContext) TakeSnapShot() interface{} {
	// take a deep copy to avoid reading & writing the same map at the same time
	k := make([]string, len(tc.DataK))
	v := make([]interface{}, len(tc.DataV))
	copy(k, tc.DataK)
	copy(v, tc.DataV)
	if tc.n == 0 {
		return &traceContext{nil, 0, k, v, nil}
	}
	last := tc.tail()
	sw := &spanWrapper{last, nil}
	return &traceContext{sw, 1, k, v, nil}
}

func GetGLocalData(key string) interface{} {
	t := getOrInitTraceContext()
	for i := 0; i < len(t.DataK); i++ {
		if t.DataK[i] == key {
			return t.DataV[i]
		}
	}
	return nil
}

func SetGLocalData(key string, value interface{}) {
	t := getOrInitTraceContext()
	t.DataK = append(t.DataK, key)
	t.DataV = append(t.DataV, value)
	if len(t.DataK) != len(t.DataV) {
		panic("DataK and DataV should have the same length")
	}
	setTraceContext(t)
}

func getOrInitTraceContext() *traceContext {
	tc := GetTraceContextFromGLS()
	if tc == nil {
		newTc := &traceContext{nil, 0, nil, nil, nil}
		setTraceContext(newTc)
		return newTc
	} else {
		return tc.(*traceContext)
	}
}

func setTraceContext(tc *traceContext) {
	SetTraceContextToGLS(tc)
}

func traceContextAddSpan(span trace.Span) {
	tc := getOrInitTraceContext()
	if !tc.add(span) {
		fmt.Println("Failed to add span to TraceContext")
	}
}

func GetTraceAndSpanId() (string, string) {
	tc := GetTraceContextFromGLS()
	if tc == nil || tc.(*traceContext).tail() == nil {
		return "", ""
	}
	ctx := tc.(*traceContext).tail().SpanContext()
	return ctx.TraceID().String(), ctx.SpanID().String()
}

func traceContextDelSpan(span trace.Span) {
	ctx := getOrInitTraceContext()
	ctx.del(span)
}

func clearTraceContext() {
	getOrInitTraceContext().clear()
}

func SpanFromGLS() trace.Span {
	gls := GetTraceContextFromGLS()
	if gls == nil {
		return nil
	}
	return gls.(*traceContext).tail()
}

func LocalRootSpanFromGLS() trace.Span {
	gls := GetTraceContextFromGLS()
	if gls == nil {
		return nil
	}
	return gls.(*traceContext).lcs
}
