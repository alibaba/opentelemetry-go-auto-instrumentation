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
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
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

const (
	KindNoop     = "noop"
	KindFile     = "file"
	KindStdout   = "stdout"
	KindZipkin   = "zipkin"
	KindOtlpGrpc = "grpc"
	KindOtlpHttp = "http"
)

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

func createSpanProcessor(ctx context.Context, endpoint, batcher string) (trace.SpanProcessor, error) {
	if verifier.IsInTest() {
		traceExporter := verifier.GetSpanExporter()
		// in test, we just send the span immediately
		simpleProcessor := trace.NewSimpleSpanProcessor(traceExporter)
		return simpleProcessor, nil
	}

	var (
		err           error
		traceExporter trace.SpanExporter
	)
	// Just support jaeger and zipkin now, more for later
	switch batcher {
	case KindZipkin:
		traceExporter, err = zipkin.New(endpoint)
	case KindOtlpGrpc:
		// Always treat trace exporter as optional component, so we use nonblock here,
		// otherwise this would slow down app start up even set a dial timeout here when
		// endpoint can not reach.
		// If the connection not dial success, the global otel ErrorHandler will catch error
		// when reporting data like other exporters.
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(endpoint),
		}
		traceExporter, err = otlptracegrpc.New(ctx, opts...)
	case KindStdout:
		traceExporter, err = stdouttrace.New()
	case KindFile:
		f, err := os.OpenFile(endpoint, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("file exporter endpoint error: %s", err.Error())
		}
		traceExporter, err = stdouttrace.New(stdouttrace.WithWriter(f))
	case KindNoop:
		traceExporter, err = tracetest.NewNoopExporter(), nil
	default:
		// Not support flexible configuration now.
		opts := []otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(endpoint),
		}
		traceExporter, err = otlptracehttp.New(ctx, opts...)
	}

	return trace.NewBatchSpanProcessor(traceExporter), err
}

func initOpenTelemetry() func() {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	reportProtocol := os.Getenv(report_protocol)
	if reportProtocol == "" {
		reportProtocol = os.Getenv(trace_report_protocol)
	}

	batchSpanProcessor, err := createSpanProcessor(ctx, endpoint, reportProtocol)
	if err != nil {
		log.Fatalf("Failed to create the OpenTelemetry trace exporter, err: %v, endpoint: %s, report_protocol: %s", err, endpoint, reportProtocol)
	}

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
