package verifier

import (
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"sort"
	"strings"
)

type SpanVerifier struct {
	verifiers            []func(tracetest.SpanStub)
	conditionalVerifiers []conditionalVerifier
}

func NewSpanVerifier() *SpanVerifier {
	return &SpanVerifier{
		verifiers: make([]func(tracetest.SpanStub), 0),
	}
}

func (sv *SpanVerifier) Verify(span tracetest.SpanStub) {
	for _, verifier := range sv.verifiers {
		verifier(span)
	}
	for _, cv := range sv.conditionalVerifiers {
		cv.Verify(span)
	}
}

func (sv *SpanVerifier) HasSpanKind(kind trace.SpanKind) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		Assert(kind == span.SpanKind, "Failed to verify span kind, expect value: %v, actual value: %v", kind, span.SpanKind)
	})
	return sv
}

func (sv *SpanVerifier) HasName(name string) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		Assert(name == span.Name, "Failed to verify span name, expect value: %v, actual value: %v", name, span.Name)
	})
	return sv
}

func (sv *SpanVerifier) HasNoParent() *SpanVerifier {
	noopSpanId := trace.SpanID{}.String()
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		Assert(span.Parent.SpanID().String() == noopSpanId, "Failed to verify parent span, expect value: %v, actual value: %v", trace.SpanID{}.String(), span.Parent.SpanID().String())
	})
	return sv
}

func (sv *SpanVerifier) HasParent(parentSpan tracetest.SpanStub) *SpanVerifier {
	pTraceId := parentSpan.SpanContext.TraceID().String()
	pSpanId := parentSpan.SpanContext.SpanID().String()
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		Assert(span.SpanContext.TraceID().String() == pTraceId, "Failed to verify parent trace id, expect value: %v, actual value: %v", pTraceId, span.SpanContext.TraceID().String())
		Assert(span.Parent.SpanID().String() == pSpanId, "Failed to verify parent span id, expect value: %v, actual value: %v", pSpanId, span.Parent.SpanID().String())
	})
	return sv
}

func (sv *SpanVerifier) HasBoolAttribute(key string, value bool) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, value)
		v := attr.AsBool()
		Assert(value == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, value, v)
	})
	return sv
}

func (sv *SpanVerifier) HasInt64Attribute(key string, value int64) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, value)
		v := attr.AsInt64()
		Assert(value == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, value, v)
	})
	return sv
}

func (sv *SpanVerifier) HasFloat64Attribute(key string, value float64) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, value)
		v := attr.AsFloat64()
		Assert(value == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, value, v)
	})
	return sv
}

func (sv *SpanVerifier) HasStringAttribute(key string, value string) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, value)
		v := attr.AsString()
		Assert(value == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, value, v)
	})
	return sv
}

func (sv *SpanVerifier) HasStringAttributeContains(key string, value string) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value contains: %v, actual value: <unset>", key, value)
		v := attr.AsString()
		Assert(strings.Contains(v, value), "Failed to verify attribute of %v, expect value contains: %v, actual value: %v", key, value, v)
	})
	return sv
}

func (sv *SpanVerifier) HasBoolSliceAttribute(key string, values []bool, ignoreOrder bool) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, values)
		vs := attr.AsBoolSlice()
		Assert(len(values) == len(vs), "Failed to verify attribute of %v, expect length: %v, actual length: %v", key, len(values), len(vs))
		if ignoreOrder {
			ac := 0
			ec := 0
			for i, v := range vs {
				if v {
					ac++
				}
				if values[i] {
					ec++
				}
			}
			Assert(ac == ec, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, values, vs)
		} else {
			for i, v := range vs {
				Assert(values[i] == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, values, vs)
			}
		}
	})
	return sv
}

func (sv *SpanVerifier) HasInt64SliceAttribute(key string, values []int64, ignoreOrder bool) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, values)
		vs := attr.AsInt64Slice()
		Assert(len(values) == len(vs), "Failed to verify attribute of %v, expect length: %v, actual length: %v", key, len(values), len(vs))
		if ignoreOrder {
			sort.Slice(values, func(i, j int) bool {
				return values[i] < values[j]
			})
			sort.Slice(vs, func(i, j int) bool {
				return vs[i] < vs[j]
			})
		}
		for i, v := range vs {
			Assert(values[i] == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, values, vs)
		}
	})
	return sv
}

func (sv *SpanVerifier) HasFloat64SliceAttribute(key string, values []float64, ignoreOrder bool) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, values)
		vs := attr.AsFloat64Slice()
		Assert(len(values) == len(vs), "Failed to verify attribute of %v, expect length: %v, actual length: %v", key, len(values), len(vs))
		if ignoreOrder {
			sort.Slice(values, func(i, j int) bool {
				return values[i] < values[j]
			})
			sort.Slice(vs, func(i, j int) bool {
				return vs[i] < vs[j]
			})
		}
		for i, v := range vs {
			Assert(values[i] == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, values, vs)
		}
	})
	return sv
}

func (sv *SpanVerifier) HasStringSliceAttribute(key string, values []string, ignoreOrder bool) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v, expect value: %v, actual value: <unset>", key, values)
		vs := attr.AsStringSlice()
		Assert(len(values) == len(vs), "Failed to verify attribute of %v, expect length: %v, actual length: %v", key, len(values), len(vs))
		if ignoreOrder {
			sort.Strings(values)
			sort.Strings(vs)
		}
		for i, v := range vs {
			Assert(values[i] == v, "Failed to verify attribute of %v, expect value: %v, actual value: %v", key, values, vs)
		}
	})
	return sv
}

func (sv *SpanVerifier) HasItemInStringSliceAttribute(key string, index int, predicate func(string) (bool, string)) *SpanVerifier {
	sv.verifiers = append(sv.verifiers, func(span tracetest.SpanStub) {
		attr := GetAttribute(span.Attributes, key)
		Assert(!IsAttributeNoop(attr), "Failed to verify attribute of %v[%v], actual value: <unset>", key, index)
		vs := attr.AsStringSlice()
		Assert(len(vs) > index, "Failed to verify attribute of %v[%v], cause index out of bound, actual length: %v", key, index, len(vs))
		result, message := predicate(vs[index])
		Assert(result, message)
	})
	return sv
}

func (sv *SpanVerifier) ConditionalVerifier(predicate func() bool, verifier *SpanVerifier) *SpanVerifier {
	sv.conditionalVerifiers = append(sv.conditionalVerifiers, conditionalVerifier{
		predicate:    predicate,
		SpanVerifier: verifier,
	})
	return sv
}

func (sv *SpanVerifier) Merge(verifier *SpanVerifier) *SpanVerifier {
	if verifier != nil {
		sv.verifiers = append(sv.verifiers, verifier.verifiers...)
		sv.conditionalVerifiers = append(sv.conditionalVerifiers, verifier.conditionalVerifiers...)
	}
	return sv
}

type conditionalVerifier struct {
	*SpanVerifier
	predicate func() bool
}

func (cv *conditionalVerifier) Verify(span tracetest.SpanStub) {
	if cv.predicate() {
		cv.SpanVerifier.Verify(span)
	}
}
