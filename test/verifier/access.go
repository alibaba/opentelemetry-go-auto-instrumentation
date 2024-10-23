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
