package instrumenter

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"testing"
	"time"
)

type testRequest struct {
}

type testResponse struct {
	status string
}

type testNameExtractor struct {
}

func (t testNameExtractor) Extract(request testRequest) string {
	return "test"
}

type testOperationListener struct {
}

func (t *testOperationListener) OnBeforeStart(parentContext context.Context, startTimestamp time.Time) context.Context {
	return context.WithValue(parentContext, "startTs", startTimestamp)
}

func (t *testOperationListener) OnBeforeEnd(ctx context.Context, startAttributes []attribute.KeyValue, startTimestamp time.Time) context.Context {
	return context.WithValue(ctx, "startAttrs", startAttributes)
}

func (t *testOperationListener) OnAfterStart(context context.Context, endTimestamp time.Time) {
	if time.Now().Sub(endTimestamp).Seconds() > 5 {
		panic("duration too long")
	}
}

func (t *testOperationListener) OnAfterEnd(context context.Context, endAttributes []attribute.KeyValue, endTimestamp time.Time) {
	if endAttributes[0].Key != "testAttribute" {
		panic("invalid attribute key")
	}
	if endAttributes[0].Value.AsString() != "testValue" {
		panic("invalid attribute value")
	}
}

type testAttributesExtractor struct {
}

func (t testAttributesExtractor) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request testRequest) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("testAttribute", "testValue"),
	}
}

func (t testAttributesExtractor) OnEnd(attributes []attribute.KeyValue, context context.Context, request testRequest, response testResponse, err error) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("testAttribute", "testValue"),
	}
}

type testContextCustomizer struct {
}

func (t testContextCustomizer) OnStart(ctx context.Context, request testRequest, startAttributes []attribute.KeyValue) context.Context {
	return context.WithValue(ctx, "test-customizer", "test-customizer")
}

type testStatusExtractor struct {
}

func (t testStatusExtractor) Extract(span trace.Span, request testRequest, response testResponse, err error) {
	if err.Error() != "abc" {
		panic(err)
	}
}

func TestInstrumenter(t *testing.T) {
	builder := Builder[testRequest, testResponse]{}
	builder.Init().
		SetSpanNameExtractor(testNameExtractor{}).
		SetSpanKindExtractor(&alwaysClientExtractor[testRequest]{}).
		AddAttributesExtractor(testAttributesExtractor{}).
		AddOperationListeners(&OperationListenerWrapper{
			listener:       &testOperationListener{},
			attrCustomizer: NoopAttrsShadower{},
		}).AddContextCustomizers(testContextCustomizer{})
	instrumenter := builder.BuildInstrumenter()
	ctx := context.Background()
	newCtx := instrumenter.Start(ctx, testRequest{})
	if newCtx.Value("test-customizer") != "test-customizer" {
		t.Fatal("key test-customizer is not expected")
	}
	if newCtx.Value("startTs") == nil {
		t.Fatal("startTs is not expected")
	}
	if newCtx.Value("startAttrs") == nil {
		t.Fatal("startAttrs is not expected")
	}
	instrumenter.End(ctx, testRequest{}, testResponse{}, errors.New("abc"))
}
