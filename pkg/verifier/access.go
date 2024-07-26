package verifier

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// In memory exporter
var spanExporter = tracetest.NewInMemoryExporter()

func GetSpanExporter() trace.SpanExporter {
	return spanExporter
}

func GetTestSpans() *tracetest.SpanStubs {
	spans := spanExporter.GetSpans()
	return &spans
}

func ResetTestSpans() {
	spanExporter.Reset()
}

func ClearSpan() {
	spanExporter.Reset()
}
