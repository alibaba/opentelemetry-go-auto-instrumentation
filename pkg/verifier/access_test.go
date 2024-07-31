package verifier

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"testing"
	"time"
)

type testSpan struct {
	// Embed the interface to implement the private method.
	sdktrace.ReadOnlySpan
	ID                   string
	name                 string
	spanContext          trace.SpanContext
	parent               trace.SpanContext
	spanKind             trace.SpanKind
	startTime            time.Time
	endTime              time.Time
	attributes           []attribute.KeyValue
	events               []sdktrace.Event
	links                []sdktrace.Link
	status               sdktrace.Status
	droppedAttributes    int
	droppedEvents        int
	droppedLinks         int
	childSpanCount       int
	resource             *resource.Resource
	instrumentationScope instrumentation.Scope
}

func (s testSpan) Name() string                     { return s.name }
func (s testSpan) SpanContext() trace.SpanContext   { return s.spanContext }
func (s testSpan) Parent() trace.SpanContext        { return s.parent }
func (s testSpan) SpanKind() trace.SpanKind         { return s.spanKind }
func (s testSpan) StartTime() time.Time             { return s.startTime }
func (s testSpan) EndTime() time.Time               { return s.endTime }
func (s testSpan) Attributes() []attribute.KeyValue { return s.attributes }
func (s testSpan) Links() []sdktrace.Link           { return s.links }
func (s testSpan) Events() []sdktrace.Event         { return s.events }
func (s testSpan) Status() sdktrace.Status          { return s.status }
func (s testSpan) DroppedAttributes() int           { return s.droppedAttributes }
func (s testSpan) DroppedLinks() int                { return s.droppedLinks }
func (s testSpan) DroppedEvents() int               { return s.droppedEvents }
func (s testSpan) ChildSpanCount() int              { return s.childSpanCount }
func (s testSpan) Resource() *resource.Resource     { return s.resource }
func (s testSpan) InstrumentationScope() instrumentation.Scope {
	return s.instrumentationScope
}

func (s testSpan) InstrumentationLibrary() instrumentation.Library {
	return s.instrumentationScope
}

func TestResetSpan(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []sdktrace.ReadOnlySpan{
		testSpan{ID: "1"},
		testSpan{ID: "2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	ResetTestSpans()
	if len(spanExporter.GetSpans()) != 0 {
		t.Fatal("expected no all the spans are cleared")
	}
}

func TestGetTestSpans(t *testing.T) {
	err := GetSpanExporter().ExportSpans(context.Background(), []sdktrace.ReadOnlySpan{
		testSpan{ID: "1"},
		testSpan{ID: "2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(spanExporter.GetSpans()) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(spanExporter.GetSpans()))
	}
}
