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

package http

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"testing"
	"time"
)

func TestHttpServerMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("test-meter")
	server, err := NewHttpServerMetric("test", meter)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	start := time.Now()
	ctx = server.OnBeforeStart(ctx, start)
	ctx = server.OnBeforeEnd(ctx, []attribute.KeyValue{}, start)
	server.OnAfterStart(ctx, start)
	server.OnAfterEnd(ctx, []attribute.KeyValue{}, time.Now())
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	if rm.ScopeMetrics[0].Metrics[0].Name != "http.server.request.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

func TestHttpClientMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("test-meter")
	client, err := NewHttpClientMetric("test", meter)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	start := time.Now()
	ctx = client.OnBeforeStart(ctx, start)
	ctx = client.OnBeforeEnd(ctx, []attribute.KeyValue{}, start)
	client.OnAfterStart(ctx, start)
	client.OnAfterEnd(ctx, []attribute.KeyValue{}, time.Now())
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	if rm.ScopeMetrics[0].Metrics[0].Name != "http.client.request.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}
