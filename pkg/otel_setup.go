//go:build ignore

package pkg

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	_ "go.opentelemetry.io/otel/sdk/trace"
)

// set the following environment variables based on https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables
// your service name: OTEL_SERVICE_NAME
// your otlp endpoint: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
const exec_name = "otel-go-auto-instrumentation"

func init() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// skip when the executable is otel-go-auto-instrumentation itself
	if strings.HasSuffix(path, exec_name) {
		return
	}
	initOpenTelemetry()
}

func newHTTPExporterAndSpanProcessor(ctx context.Context) (trace.SpanExporter, trace.SpanProcessor) {
	if verifier.IsInTest() {
		traceExporter := verifier.GetSpanExporter()
		// in test, we just send the span immediately
		simpleProcessor := trace.NewSimpleSpanProcessor(traceExporter)
		return traceExporter, simpleProcessor
	} else {
		traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
		if err != nil {
			log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
		}
		batchSpanProcessor := trace.NewBatchSpanProcessor(traceExporter)
		return traceExporter, batchSpanProcessor
	}
}

func initOpenTelemetry() func() {
	ctx := context.Background()

	var traceExporter trace.SpanExporter
	var batchSpanProcessor trace.SpanProcessor

	traceExporter, batchSpanProcessor = newHTTPExporterAndSpanProcessor(ctx)

	// TODO: add sampler
	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(batchSpanProcessor))

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExporter.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
