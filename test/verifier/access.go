// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package verifier

import (
	"context"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// In memory span exporter
var spanExporter = tracetest.NewInMemoryExporter()

var ManualReader = metric.NewManualReader()

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

func GetTestMetrics() (metricdata.ResourceMetrics, error) {
	var tmp, result metricdata.ResourceMetrics
	err := ManualReader.Collect(context.Background(), &tmp)
	if err != nil {
		return metricdata.ResourceMetrics{}, err
	}
	result = deepCopyMetric(tmp)
	return result, nil
}
