//go:build ignore

package pkg

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
	"strings"
	"time"
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

func newHTTPExporterAndSpanProcessor(ctx context.Context) (*otlptrace.Exporter, trace.SpanProcessor) {

	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())

	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
	}

	batchSpanProcessor := trace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor
}

func initOpenTelemetry() func() {
	ctx := context.Background()

	var traceExporter *otlptrace.Exporter
	var batchSpanProcessor trace.SpanProcessor

	traceExporter, batchSpanProcessor = newHTTPExporterAndSpanProcessor(ctx)

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
