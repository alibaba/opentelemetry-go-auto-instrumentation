package verifier

import (
	"context"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
	"testing"
)

func TestWaitAndAssertTracesOneTrace(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []trace.ReadOnlySpan{
		testSpan{ID: "1", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "2", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "3", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "4", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ResetTestSpans()
	WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) != 1 {
			t.Fatalf("expecting 1 traces but got %d", len(stubs))
		}
	})
}

func TestWaitAndAssertTracesMultipleTrace(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []trace.ReadOnlySpan{
		testSpan{ID: "1", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "2", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x020},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "3", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x030},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "4", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x040},
			SpanID:  oteltrace.SpanID{0x01},
		})},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ResetTestSpans()
	WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) != 4 {
			t.Fatalf("expecting 4 traces but got %d", len(stubs))
		}
	})
}

func TestWaitAndAssertTraceLink(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []trace.ReadOnlySpan{
		testSpan{ID: "1", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "2", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x02},
		}), parent: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ResetTestSpans()
	WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) != 1 {
			t.Fatalf("expecting 1 traces but got %d", len(stubs))
		}
		if stubs[0].Snapshots()[1].Parent().SpanID() != stubs[0].Snapshots()[0].SpanContext().SpanID() {
			t.Fatalf("expecting parent span id to be equal")
		}
	})
}
