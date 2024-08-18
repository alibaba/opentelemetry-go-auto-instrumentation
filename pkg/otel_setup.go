// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:build ignore

package pkg

import (
	"context"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	_ "go.opentelemetry.io/otel/sdk/trace"
)

// set the following environment variables based on https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables
// your service name: OTEL_SERVICE_NAME
// your otlp endpoint: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
// your otlp header: OTEL_EXPORTER_OTLP_HEADERS
const exec_name = "otelbuild"
const report_protocol = "OTEL_EXPORTER_OTLP_PROTOCOL"
const trace_report_protocol = "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"

func init() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// skip when the executable is otelbuild itself
	if strings.HasSuffix(path, exec_name) {
		return
	}
	initOpenTelemetry()
}

func newSpanProcessor(ctx context.Context) trace.SpanProcessor {
	if verifier.IsInTest() {
		traceExporter := verifier.GetSpanExporter()
		// in test, we just send the span immediately
		simpleProcessor := trace.NewSimpleSpanProcessor(traceExporter)
		return simpleProcessor
	} else {
		var traceExporter trace.SpanExporter
		var err error
		if os.Getenv(report_protocol) == "grpc" || os.Getenv(trace_report_protocol) == "grpc" {
			traceExporter, err = otlptrace.New(ctx, otlptracegrpc.NewClient())
			if err != nil {
				log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
			}
		} else {
			traceExporter, err = otlptrace.New(ctx, otlptracehttp.NewClient())
			if err != nil {
				log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
			}
		}
		batchSpanProcessor := trace.NewBatchSpanProcessor(traceExporter)
		return batchSpanProcessor
	}
}

func initOpenTelemetry() func() {
	ctx := context.Background()

	var batchSpanProcessor trace.SpanProcessor

	batchSpanProcessor = newSpanProcessor(ctx)

	var traceProvider *trace.TracerProvider
	if batchSpanProcessor != nil {
		traceProvider = trace.NewTracerProvider(
			trace.WithSpanProcessor(batchSpanProcessor))
	} else {
		traceProvider = trace.NewTracerProvider()
	}

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
