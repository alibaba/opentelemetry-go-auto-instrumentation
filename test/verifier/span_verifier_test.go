package verifier

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

var (
	// 80f198ee56343ba864fe8b2a57d3eff7
	ValidTraceID = trace.TraceID([16]byte{128, 241, 152, 238, 86, 52, 59, 168, 100, 254, 139, 42, 87, 211, 239, 247})
	// aa0ba902b700f067
	ValidSpanID = trace.SpanID([8]byte{170, 11, 169, 2, 183, 0, 240, 103})
	// 00f067aa0ba902b7
	ValidChildSpanID = trace.SpanID([8]byte{0, 240, 103, 170, 11, 169, 2, 183})
)

func TestSpanVerifier_HasSpanKind(t *testing.T) {
	span := newSpanStub(nil, "testSpan", trace.SpanKindClient)

	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightSpanKind", func() {
		sv.HasSpanKind(trace.SpanKindClient).Verify(span)
	}, true)

	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongSpanKind", func() {
		sv.HasSpanKind(trace.SpanKindInternal).Verify(span)
	}, false)
}

func TestSpanVerifier_HasName(t *testing.T) {
	span := newSpanStub(nil, "testSpan", trace.SpanKindInternal)

	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightSpanName", func() {
		sv.HasName("testSpan").Verify(span)
	}, true)

	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongSpanName", func() {
		sv.HasName("wrongSpan").Verify(span)
	}, false)
}

func TestSpanVerifier_HasNoParent(t *testing.T) {
	// HasNoParent
	sv := NewSpanVerifier()
	span := newSpanStub(nil, "testSpan", trace.SpanKindInternal)

	verifyVerifier(t, "TestNoParent", func() {
		sv.HasNoParent().Verify(span)
	}, true)

	// HasParent
	sv = NewSpanVerifier()
	parentSpan := newSpanStub(nil, "parentSpan", trace.SpanKindInternal)
	span = newSpanStub(&parentSpan, "testSpan", trace.SpanKindInternal)

	verifyVerifier(t, "TestHasParent", func() {
		sv.HasNoParent().Verify(span)
	}, false)
}

func TestSpanVerifier_HasParent(t *testing.T) {
	// HasParent
	sv := NewSpanVerifier()
	parentSpan := newSpanStub(nil, "parentSpan", trace.SpanKindInternal)
	span := newSpanStub(&parentSpan, "testSpan", trace.SpanKindInternal)
	verifyVerifier(t, "TestHasParent", func() {
		sv.HasParent(parentSpan).Verify(span)
	}, true)

	// HasNoParent
	sv = NewSpanVerifier()
	span = newSpanStub(nil, "testSpan", trace.SpanKindInternal)
	verifyVerifier(t, "TestNoParent", func() {
		sv.HasParent(parentSpan).Verify(span)
	}, false)
}

func TestSpanVerifier_HasBoolAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.Bool("boolKeyTrue", true),
		attribute.Bool("boolKeyFalse", false),
	)

	// RightBoolKey
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightBoolKey", func() {
		sv.HasBoolAttribute("boolKeyTrue", true).
			HasBoolAttribute("boolKeyFalse", false).
			Verify(span)
	}, true)

	// WrongBoolKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongBoolKey", func() {
		sv.HasBoolAttribute("boolKeyTrue", false).
			HasBoolAttribute("boolKeyFalse", true).
			Verify(span)
	}, false)

	// NilBoolKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilBoolKey", func() {
		sv.HasBoolAttribute("boolKeyOther", false).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasInt64Attribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.Int("IntKey", 42),
		attribute.Int("IntKey0", 0),
		attribute.Int64("Int64Key", 42),
		attribute.Int64("Int64Key0", 0),
	)

	// RightInt64Key
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightInt64Key", func() {
		sv.HasInt64Attribute("IntKey", 42).
			HasInt64Attribute("IntKey0", 0).
			HasInt64Attribute("Int64Key", 42).
			HasInt64Attribute("Int64Key0", 0).
			Verify(span)
	}, true)

	// WrongInt64Key
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongInt64Key", func() {
		sv.HasInt64Attribute("IntKey", 0).
			HasInt64Attribute("IntKey0", 42).
			HasInt64Attribute("Int64Key", 0).
			HasInt64Attribute("Int64Key0", 42).
			Verify(span)
	}, false)

	// NilInt64Key
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilInt64Key", func() {
		sv.HasInt64Attribute("Int64KeyOther", 0).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasFloat64Attribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.Float64("Float64Key", 42.42),
		attribute.Float64("Float64Key0", 0),
	)

	// RightFloat64Key
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightFloat64Key", func() {
		sv.HasFloat64Attribute("Float64Key", 42.42).
			HasFloat64Attribute("Float64Key0", 0).
			Verify(span)
	}, true)

	// WrongFloat64Key
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongFloat64Key", func() {
		sv.HasFloat64Attribute("Float64Key", 0).
			HasFloat64Attribute("Float64Key0", 42.42).
			Verify(span)
	}, false)

	// NilFloat64Key
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilFloat64Key", func() {
		sv.HasFloat64Attribute("Float64KeyOther", 0).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasStringAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.String("StringKey", "testValue"),
		attribute.String("StringKey0", ""),
	)

	// RightStringKey
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightStringKey", func() {
		sv.HasStringAttribute("StringKey", "testValue").
			HasStringAttribute("StringKey0", "").
			Verify(span)
	}, true)

	// WrongStringKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongStringKey", func() {
		sv.HasStringAttribute("StringKey", "").
			HasStringAttribute("StringKey0", "testValue").
			Verify(span)
	}, false)

	// NilStringKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilStringKey", func() {
		sv.HasStringAttribute("StringKeyOther", "").
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasStringAttributeContains(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.String("StringKey", "mouse|cat|dog"),
		attribute.String("StringKey0", ""),
	)

	// RightStringKey
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightStringKey", func() {
		sv.HasStringAttributeContains("StringKey", "cat").
			HasStringAttributeContains("StringKey", "mouse").
			HasStringAttribute("StringKey0", "").
			Verify(span)
	}, true)

	// WrongStringKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongStringKey", func() {
		sv.HasStringAttribute("StringKey", "sun").
			HasStringAttribute("StringKey0", "sun").
			Verify(span)
	}, false)

	// NilStringKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilStringKey", func() {
		sv.HasStringAttribute("StringKeyOther", "sun").
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasBoolSliceAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.BoolSlice("BoolSliceKey", []bool{true, true, false, true}),
		attribute.BoolSlice("BoolSliceKey0", []bool{}),
	)

	// RightBoolSliceKey with order
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightBoolSliceKeyWithOrder", func() {
		sv.HasBoolSliceAttribute("BoolSliceKey", []bool{true, true, false, true}, false).
			HasBoolSliceAttribute("BoolSliceKey0", []bool{}, false).
			Verify(span)
	}, true)

	// RightBoolSliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestRightBoolSliceKeyWithoutOrder", func() {
		sv.HasBoolSliceAttribute("BoolSliceKey", []bool{true, true, true, false}, true).
			HasBoolSliceAttribute("BoolSliceKey0", []bool{}, true).
			Verify(span)
	}, true)

	// WrongBoolSliceKey with order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongBoolSliceKeyWithOrder", func() {
		sv.HasBoolSliceAttribute("BoolSliceKey", []bool{}, false).
			HasBoolSliceAttribute("BoolSliceKey0", []bool{true, true, false, true}, false).
			Verify(span)
	}, false)

	// WrongBoolSliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongBoolSliceKeyWithoutOrder", func() {
		sv.HasBoolSliceAttribute("BoolSliceKey", []bool{}, true).
			HasBoolSliceAttribute("BoolSliceKey0", []bool{true, true, true, false}, true).
			Verify(span)
	}, false)

	// NilBoolSliceKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilBoolSliceKey", func() {
		sv.HasBoolSliceAttribute("BoolSliceKeyOther", []bool{}, true).
			HasBoolSliceAttribute("BoolSliceKeyOther2", []bool{}, false).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasInt64SliceAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.IntSlice("IntSliceKey", []int{42, 12, 4242}),
		attribute.IntSlice("IntSliceKey0", []int{}),
		attribute.Int64Slice("Int64SliceKey", []int64{42, 12, 4242}),
		attribute.Int64Slice("Int64SliceKey0", []int64{}),
	)

	// RightInt64SliceKey with order
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightInt64SliceKeyWithOrder", func() {
		sv.HasInt64SliceAttribute("IntSliceKey", []int64{42, 12, 4242}, false).
			HasInt64SliceAttribute("IntSliceKey0", []int64{}, false).
			HasInt64SliceAttribute("Int64SliceKey", []int64{42, 12, 4242}, false).
			HasInt64SliceAttribute("Int64SliceKey0", []int64{}, false).
			Verify(span)
	}, true)

	// RightInt64SliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestRightInt64SliceKeyWithoutOrder", func() {
		sv.HasInt64SliceAttribute("IntSliceKey", []int64{4242, 12, 42}, true).
			HasInt64SliceAttribute("IntSliceKey0", []int64{}, true).
			HasInt64SliceAttribute("Int64SliceKey", []int64{4242, 12, 42}, true).
			HasInt64SliceAttribute("Int64SliceKey0", []int64{}, true).
			Verify(span)
	}, true)

	// WrongInt64SliceKey with order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongInt64SliceKeyWithOrder", func() {
		sv.HasInt64SliceAttribute("IntSliceKey", []int64{}, false).
			HasInt64SliceAttribute("IntSliceKey0", []int64{42, 12, 4242}, false).
			HasInt64SliceAttribute("Int64SliceKey", []int64{}, false).
			HasInt64SliceAttribute("Int64SliceKey0", []int64{42, 12, 4242}, false).
			Verify(span)
	}, false)

	// WrongInt64SliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongInt64SliceKeyWithOrder", func() {
		sv.HasInt64SliceAttribute("IntSliceKey", []int64{}, true).
			HasInt64SliceAttribute("IntSliceKey0", []int64{4242, 12, 42}, true).
			HasInt64SliceAttribute("Int64SliceKey", []int64{}, true).
			HasInt64SliceAttribute("Int64SliceKey0", []int64{4242, 12, 42}, true).
			Verify(span)
	}, false)

	// NilInt64SliceKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilInt64SliceKey", func() {
		sv.HasInt64SliceAttribute("Int64SliceKeyOther", []int64{}, true).
			HasInt64SliceAttribute("Int64SliceKeyOther2", []int64{}, false).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasFloat64SliceAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.Float64Slice("Float64SliceKey", []float64{12.1, 42.1, -1.2, -98.3}),
		attribute.Float64Slice("Float64SliceKey0", []float64{}),
	)

	// RightFloat64SliceKey with order
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightFloat64SliceKeyWithOrder", func() {
		sv.HasFloat64SliceAttribute("Float64SliceKey", []float64{12.1, 42.1, -1.2, -98.3}, false).
			HasFloat64SliceAttribute("Float64SliceKey0", []float64{}, false).
			Verify(span)
	}, true)

	// RightFloat64SliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestRightFloat64SliceKeyWithoutOrder", func() {
		sv.HasFloat64SliceAttribute("Float64SliceKey", []float64{-98.3, 42.1, 12.1, -1.2}, true).
			HasFloat64SliceAttribute("Float64SliceKey0", []float64{}, true).
			Verify(span)
	}, true)

	// WrongFloat64SliceKey with order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongFloat64SliceKeyWithOrder", func() {
		sv.HasFloat64SliceAttribute("Float64SliceKey", []float64{}, false).
			HasFloat64SliceAttribute("Float64SliceKey0", []float64{12.1, 42.1, -1.2, -98.3}, false).
			Verify(span)
	}, false)

	// WrongFloat64SliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongFloat64SliceKeyWithoutOrder", func() {
		sv.HasFloat64SliceAttribute("Float64SliceKey", []float64{}, true).
			HasFloat64SliceAttribute("Float64SliceKey0", []float64{-98.3, 42.1, 12.1, -1.2}, true).
			Verify(span)
	}, false)

	// NilFloat64SliceKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilFloat64SliceKey", func() {
		sv.HasFloat64SliceAttribute("Float64SliceKeyOther", []float64{}, true).
			HasFloat64SliceAttribute("Float64SliceKeyOther2", []float64{}, false).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasStringSliceAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.StringSlice("StringSliceKey", []string{"w3c", "b3", "jaeger", "opentracing"}),
		attribute.StringSlice("StringSliceKey0", []string{}),
	)

	// RightStringSliceKey with order
	sv := NewSpanVerifier()
	verifyVerifier(t, "TestRightStringSliceKeyWithOrder", func() {
		sv.HasStringSliceAttribute("StringSliceKey", []string{"w3c", "b3", "jaeger", "opentracing"}, false).
			HasStringSliceAttribute("StringSliceKey0", []string{}, false).
			Verify(span)
	}, true)

	// RightStringSliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestRightStringSliceKeyWithoutOrder", func() {
		sv.HasStringSliceAttribute("StringSliceKey", []string{"w3c", "opentracing", "jaeger", "b3"}, true).
			HasStringSliceAttribute("StringSliceKey0", []string{}, true).
			Verify(span)
	}, true)

	// WrongStringSliceKey with order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongStringSliceKeyWithOrder", func() {
		sv.HasStringSliceAttribute("StringSliceKey", []string{}, false).
			HasStringSliceAttribute("StringSliceKey0", []string{"w3c", "b3", "jaeger", "opentracing"}, false).
			Verify(span)
	}, false)

	// WrongStringSliceKey without order
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestWrongStringSliceKeyWithoutOrder", func() {
		sv.HasStringSliceAttribute("StringSliceKey", []string{}, true).
			HasStringSliceAttribute("StringSliceKey0", []string{"w3c", "opentracing", "jaeger", "b3"}, true).
			Verify(span)
	}, false)

	// NilStringSliceKey
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestNilStringSliceKey", func() {
		sv.HasStringSliceAttribute("StringSliceKeyOther", []string{}, true).
			HasStringSliceAttribute("StringSliceKeyOther2", []string{}, false).
			Verify(span)
	}, false)
}

func TestSpanVerifier_HasItemInStringSliceAttribute(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.StringSlice("StringSliceKey", []string{"w3c", "b3", "jaeger", "opentracing"}),
		attribute.StringSlice("StringSliceKey0", []string{}),
	)

	sv := NewSpanVerifier()
	verifyVerifier(t, "TestSucc", func() {
		sv.HasItemInStringSliceAttribute("StringSliceKey", 0, func(s string) (bool, string) {
			return s == "w3c", ""
		}).
			HasItemInStringSliceAttribute("StringSliceKey", 2, func(s string) (bool, string) {
				return s == "jaeger", ""
			}).
			Verify(span)
	}, true)

	sv = NewSpanVerifier()
	verifyVerifier(t, "TestFail", func() {
		sv.HasItemInStringSliceAttribute("StringSliceKey", 2, func(s string) (bool, string) {
			return s == "w3c", ""
		}).
			HasItemInStringSliceAttribute("StringSliceKey", 0, func(s string) (bool, string) {
				return s == "jaeger", ""
			}).
			Verify(span)
	}, false)

	sv = NewSpanVerifier()
	verifyVerifier(t, "TestOutOfBound", func() {
		sv.HasItemInStringSliceAttribute("StringSliceKey0", 0, func(s string) (bool, string) {
			return s == "w3c", ""
		}).
			Verify(span)
	}, false)
}

func TestSpanVerifier_ConditionalVerifier(t *testing.T) {
	var verifyDetail = true
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.String("StringKey", "normal"),
		attribute.String("StringDetailKey", "detail"),
	)

	sv := NewSpanVerifier()
	verifyVerifier(t, "TestConditionalVerifier", func() {
		sv.HasStringAttribute("StringKey", "normal").
			ConditionalVerifier(func() bool {
				return verifyDetail
			}, NewSpanVerifier().HasStringAttribute("StringDetailKey", "detail")).
			Verify(span)
	}, true)

	verifyDetail = false
	sv = NewSpanVerifier()
	verifyVerifier(t, "TestConditionalVerifierIgnore", func() {
		sv.HasStringAttribute("StringKey", "normal").
			ConditionalVerifier(func() bool {
				return verifyDetail
			}, NewSpanVerifier().HasStringAttribute("StringDetailKey2", "detail")).
			Verify(span)
	}, true)
}

func TestSpanVerifier_Merge(t *testing.T) {
	span := newSpanStub(
		nil,
		"testSpan",
		trace.SpanKindInternal,
		attribute.String("StringKey", "normal"),
		attribute.String("StringDetailKey", "detail"),
	)

	sv := NewSpanVerifier()
	verifyVerifier(t, "TestMergeSuccVerifier", func() {
		sv.HasStringAttribute("StringKey", "normal").
			Merge(NewSpanVerifier().HasStringAttribute("StringDetailKey", "detail")).
			Verify(span)
	}, true)

	sv = NewSpanVerifier()
	verifyVerifier(t, "TestMergeFailVerifier", func() {
		sv.HasStringAttribute("StringKey", "normal").
			Merge(NewSpanVerifier().HasStringAttribute("StringDetailKey2", "detail")).
			Verify(span)
	}, false)
}

func newSpanStub(pSpan *tracetest.SpanStub, name string, kind trace.SpanKind, attrs ...attribute.KeyValue) tracetest.SpanStub {
	if pSpan == nil {
		return tracetest.SpanStub{
			SpanKind:   kind,
			Name:       name,
			Attributes: attrs,
			Parent:     trace.SpanContext{},
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    ValidTraceID,
				SpanID:     ValidSpanID,
				TraceFlags: trace.FlagsSampled,
			}),
		}
	} else {
		return tracetest.SpanStub{
			SpanKind:   kind,
			Name:       name,
			Attributes: attrs,
			Parent:     pSpan.SpanContext,
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    pSpan.SpanContext.TraceID(),
				SpanID:     ValidChildSpanID,
				TraceFlags: pSpan.SpanContext.TraceFlags(),
				TraceState: pSpan.SpanContext.TraceState(),
			}),
		}
	}
}

func verifyVerifier(t *testing.T, testName string, testFunc func(), shouldPass bool) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		if shouldPass {
			Assert(err == nil, "Failed to run %v:%v, expect result: SUCCESS, actual: FAILED, cause: %v", t.Name(), testName, err)
		} else {
			Assert(err != nil, "Failed to run %v:%v, expect result: FAILED, actual: SUCCESS", t.Name(), testName)
		}
	}()
	testFunc()
}
