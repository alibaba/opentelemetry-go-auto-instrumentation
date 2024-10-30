// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"context"
	"embed"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/core/meter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/http"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	_ "go.opentelemetry.io/otel/sdk/trace"
)

// set the following environment variables based on https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables
// your service name: OTEL_SERVICE_NAME
// your otlp endpoint: OTEL_EXPORTER_OTLP_ENDPOINT OTEL_EXPORTER_OTLP_TRACES_ENDPOINT OTEL_EXPORTER_OTLP_METRICS_ENDPOINT OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// your otlp header: OTEL_EXPORTER_OTLP_HEADERS
const exec_name = "otelbuild"
const report_protocol = "OTEL_EXPORTER_OTLP_PROTOCOL"
const trace_report_protocol = "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"

//go:embed core inst-api inst-api-semconv
var embedDir embed.FS

func ExportPkgDirList() []string {
	return []string{"core", "inst-api", "inst-api-semconv"}
}
func ExportPkgFS() embed.FS {
	return embedDir
}

//go:embed otel_setup.go
var otelSetupSDKTemplate string

func ExportOtelSetupSDKTemplate() string {
	return otelSetupSDKTemplate
}

//go:embed api/api.go
var apiSnippet string

func ExportAPISnippet() string { return apiSnippet }

//go:embed data/default.json
var defaultRuleJson string

func ExportDefaultRuleJson() string { return defaultRuleJson }

// This is simply a cache for the embedded rules directory. Technically, now any
// rules can by fetched from either network or local file system. But we want
// to keep the rules in tool binary, in this way, we dont need to fetch them from
// network, which is a very slow process.
//
//go:embed rules
var ruleCache embed.FS

func ExportRuleCache() embed.FS {
	return ruleCache
}

func init() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// skip when the executable is otelbuild itself
	if strings.HasSuffix(path, exec_name) {
		return
	}
	if err = initOpenTelemetry(); err != nil {
		log.Fatalf("%s: %v", "Failed to initialize opentelemetry resource", err)
	}
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

func initOpenTelemetry() error {
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
	return initMetrics()
}

func initMetrics() error {
	ctx := context.Background()
	var mp *metric.MeterProvider
	// TODO: abstract the if-else
	if verifier.IsInTest() {
		mp = metric.NewMeterProvider(
			metric.WithReader(verifier.ManualReader),
		)
	} else {
		if os.Getenv(report_protocol) == "grpc" || os.Getenv(trace_report_protocol) == "grpc" {
			exporter, err := otlpmetricgrpc.New(ctx)
			if err != nil {
				log.Fatalf("new otlp metric grpc exporter failed: %v", err)
			}
			mp = metric.NewMeterProvider(
				metric.WithReader(metric.NewPeriodicReader(exporter)),
			)
		} else {
			exporter, err := otlpmetrichttp.New(ctx)
			if err != nil {
				log.Fatalf("new otlp metric http exporter failed: %v", err)
			}
			mp = metric.NewMeterProvider(
				metric.WithReader(metric.NewPeriodicReader(exporter)),
			)
		}
	}
	if mp == nil {
		return errors.New("No MeterProvider is provided")
	}
	otel.SetMeterProvider(mp)
	m := mp.Meter("opentelemetry-global-meter")
	meter.SetMeter(m)
	// init http metrics
	http.InitHttpMetrics(m)
	// DefaultMinimumReadMemStatsInterval is 15 second
	return runtime.Start(runtime.WithMeterProvider(mp))
}
