//go:build ignore

package trace

import (
	"fmt"
	trace "go.opentelemetry.io/otel/trace"
)

const maxSpans = 300

type traceContext struct {
	sw   *spanWrapper
	n    int
	Data map[string]interface{}
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
	tc.sw = wrapper
	tc.n++
	return true
}

func (tc *traceContext) tail() trace.Span {
	if tc.n == 0 {
		return nil
	} else {
		return tc.sw.span
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
	tc.Data = nil
	SetBaggageContainerToGLS(nil)
}

func (tc *traceContext) TakeSnapShot() interface{} {
	// take a deep copy to avoid reading & writing the same map at the same time
	var dataCopy = make(map[string]interface{})
	for key, value := range tc.Data {
		dataCopy[key] = value
	}
	if tc.n == 0 {
		return &traceContext{nil, 0, dataCopy}
	}
	last := tc.tail()
	sw := &spanWrapper{last, nil}
	return &traceContext{sw, 1, dataCopy}
}

func GetGLocalData(key string) interface{} {
	t := getOrInitTraceContext()
	r := t.Data[key]
	return r
}

func SetGLocalData(key string, value interface{}) {
	t := getOrInitTraceContext()
	if t.Data == nil {
		t.Data = make(map[string]interface{})
	}
	t.Data[key] = value
	setTraceContext(t)
}

func getOrInitTraceContext() *traceContext {
	tc := GetTraceContextFromGLS()
	if tc == nil {
		newTc := &traceContext{nil, 0, nil}
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
