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

package rpc

import (
	"context"
	"errors"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

type testRequest struct {
}

type testResponse struct {
}

type grpcAttrsGetter struct {
}

func (h grpcAttrsGetter) GetSystem(request testRequest) string {
	return "grpc"
}

func (h grpcAttrsGetter) GetService(request testRequest) string {
	return "TestService"
}

func (h grpcAttrsGetter) GetMethod(request testRequest) string {
	return "TestMethod"
}

func (h grpcAttrsGetter) GetServerAddress(request testRequest) string {
	return "localhost:8080"
}

func TestRpcServerMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("test-meter")
	server, err := newRpcServerMetric("test", meter)
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
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.server.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

func TestRpcClientMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("test-meter")
	client, err := newRpcClientMetric("test", meter)
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
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.client.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

// Test gRPC status code in metrics with successful call
func TestGrpcMetricsWithStatusCodeSuccess(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("grpc-service"),
		semconv.ServiceVersion("v1.0.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("grpc-test-meter")
	InitRpcMetrics(meter)
	
	// Create RPC extractor for gRPC
	rpcExtractor := ClientRpcAttrsExtractor[testRequest, testResponse, grpcAttrsGetter]{}
	
	// Create client metric
	client := RpcClientMetrics("grpc.client")
	
	ctx := context.Background()
	start := time.Now()
	
	// Start phase
	startAttrs, ctx := rpcExtractor.OnStart([]attribute.KeyValue{}, ctx, testRequest{})
	ctx = client.OnBeforeStart(ctx, start)
	ctx = client.OnBeforeEnd(ctx, startAttrs, start)
	client.OnAfterStart(ctx, start)
	
	// End phase - successful call (no error)
	endAttrs, _ := rpcExtractor.OnEnd([]attribute.KeyValue{}, ctx, testRequest{}, testResponse{}, nil)
	client.OnAfterEnd(ctx, endAttrs, time.Now())
	
	// Collect metrics
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	
	// Verify metric name
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.client.duration" {
		t.Fatalf("Expected metric name rpc.client.duration, got %s", rm.ScopeMetrics[0].Metrics[0].Name)
	}
	
	// Verify gRPC status code attribute is present in the metrics
	histogram := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
	attrs := histogram.DataPoints[0].Attributes
	
	// Check for gRPC status code attribute (should be 0 for success)
	statusCodeFound := false
	for _, attr := range attrs.ToSlice() {
		if attr.Key == semconv.RPCGRPCStatusCodeKey {
			statusCodeFound = true
			if attr.Value.AsInt64() != 0 {
				t.Fatalf("Expected gRPC status code 0 (OK), got %d", attr.Value.AsInt64())
			}
			break
		}
	}
	
	if !statusCodeFound {
		t.Fatalf("gRPC status code attribute not found in metrics")
	}
}

// Test gRPC status code in metrics with error call
func TestGrpcMetricsWithStatusCodeError(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("grpc-service"),
		semconv.ServiceVersion("v1.0.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("grpc-test-meter")
	InitRpcMetrics(meter)
	
	// Create RPC extractor for gRPC
	rpcExtractor := ServerRpcAttrsExtractor[testRequest, testResponse, grpcAttrsGetter]{}
	
	// Create server metric
	server := RpcServerMetrics("grpc.server")
	
	ctx := context.Background()
	start := time.Now()
	
	// Start phase
	startAttrs, ctx := rpcExtractor.OnStart([]attribute.KeyValue{}, ctx, testRequest{})
	ctx = server.OnBeforeStart(ctx, start)
	ctx = server.OnBeforeEnd(ctx, startAttrs, start)
	server.OnAfterStart(ctx, start)
	
	// End phase - error call
	grpcErr := status.Error(codes.NotFound, "resource not found")
	endAttrs, _ := rpcExtractor.OnEnd([]attribute.KeyValue{}, ctx, testRequest{}, testResponse{}, grpcErr)
	server.OnAfterEnd(ctx, endAttrs, time.Now())
	
	// Collect metrics
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	
	// Verify metric name
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.server.duration" {
		t.Fatalf("Expected metric name rpc.server.duration, got %s", rm.ScopeMetrics[0].Metrics[0].Name)
	}
	
	// Verify gRPC status code attribute is present in the metrics
	histogram := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
	attrs := histogram.DataPoints[0].Attributes
	
	// Check for gRPC status code attribute (should be 5 for NotFound)
	statusCodeFound := false
	for _, attr := range attrs.ToSlice() {
		if attr.Key == semconv.RPCGRPCStatusCodeKey {
			statusCodeFound = true
			if attr.Value.AsInt64() != int64(codes.NotFound) {
				t.Fatalf("Expected gRPC status code %d (NotFound), got %d", int64(codes.NotFound), attr.Value.AsInt64())
			}
			break
		}
	}
	
	if !statusCodeFound {
		t.Fatalf("gRPC status code attribute not found in metrics")
	}
}

// Test non-gRPC system should not have status code in metrics
func TestNonGrpcSystemNoStatusCodeInMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("non-grpc-service"),
		semconv.ServiceVersion("v1.0.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	meter := mp.Meter("non-grpc-test-meter")
	InitRpcMetrics(meter)
	
	// Create non-gRPC attrs getter
	type nonGrpcAttrsGetter struct{}
	func (h nonGrpcAttrsGetter) GetSystem(request testRequest) string { return "other" }
	func (h nonGrpcAttrsGetter) GetService(request testRequest) string { return "TestService" }
	func (h nonGrpcAttrsGetter) GetMethod(request testRequest) string { return "TestMethod" }
	func (h nonGrpcAttrsGetter) GetServerAddress(request testRequest) string { return "localhost:8080" }
	
	// Create RPC extractor for non-gRPC system
	rpcExtractor := ClientRpcAttrsExtractor[testRequest, testResponse, nonGrpcAttrsGetter]{}
	
	// Create client metric
	client := RpcClientMetrics("other.client")
	
	ctx := context.Background()
	start := time.Now()
	
	// Start phase
	startAttrs, ctx := rpcExtractor.OnStart([]attribute.KeyValue{}, ctx, testRequest{})
	ctx = client.OnBeforeStart(ctx, start)
	ctx = client.OnBeforeEnd(ctx, startAttrs, start)
	client.OnAfterStart(ctx, start)
	
	// End phase - with any error
	endAttrs, _ := rpcExtractor.OnEnd([]attribute.KeyValue{}, ctx, testRequest{}, testResponse{}, errors.New("some error"))
	client.OnAfterEnd(ctx, endAttrs, time.Now())
	
	// Collect metrics
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	
	// Verify metric name
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.client.duration" {
		t.Fatalf("Expected metric name rpc.client.duration, got %s", rm.ScopeMetrics[0].Metrics[0].Name)
	}
	
	// Verify gRPC status code attribute is NOT present in the metrics for non-gRPC system
	histogram := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
	attrs := histogram.DataPoints[0].Attributes
	
	// Check that gRPC status code attribute is NOT present
	for _, attr := range attrs.ToSlice() {
		if attr.Key == semconv.RPCGRPCStatusCodeKey {
			t.Fatalf("gRPC status code attribute should not be present for non-gRPC system")
		}
	}
}

func TestRpcMetricAttributesShadower(t *testing.T) {
	attrs := make([]attribute.KeyValue, 0)
	attrs = append(attrs, attribute.KeyValue{
		Key:   semconv.RPCMethodKey,
		Value: attribute.StringValue("method"),
	}, attribute.KeyValue{
		Key:   "unknown",
		Value: attribute.Value{},
	}, attribute.KeyValue{
		Key:   semconv.RPCServiceKey,
		Value: attribute.StringValue("rpc"),
	}, attribute.KeyValue{
		Key:   semconv.RPCSystemKey,
		Value: attribute.StringValue("abc"),
	})
	n, attrs := utils.Shadow(attrs, rpcMetricsConv)
	if n != 3 {
		panic("wrong shadow array")
	}
	if attrs[3].Key != "unknown" {
		panic("unknown should be the last attribute")
	}
}

func TestLazyRpcServerMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	m := mp.Meter("test-meter")
	InitRpcMetrics(m)
	server := RpcServerMetrics("net.rpc.server")
	ctx := context.Background()
	start := time.Now()
	ctx = server.OnBeforeStart(ctx, start)
	ctx = server.OnBeforeEnd(ctx, []attribute.KeyValue{}, start)
	server.OnAfterStart(ctx, start)
	server.OnAfterEnd(ctx, []attribute.KeyValue{}, time.Now())
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.server.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

func TestLazyRpcClientMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	m := mp.Meter("test-meter")
	InitRpcMetrics(m)
	client := RpcClientMetrics("net.rpc.client")
	ctx := context.Background()
	start := time.Now()
	ctx = client.OnBeforeStart(ctx, start)
	ctx = client.OnBeforeEnd(ctx, []attribute.KeyValue{}, start)
	client.OnAfterStart(ctx, start)
	client.OnAfterEnd(ctx, []attribute.KeyValue{}, time.Now())
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.client.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

func TestGlobalRpcServerMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	m := mp.Meter("test-meter")
	InitRpcMetrics(m)
	server := RpcServerMetrics("net.rpc.server")
	ctx := context.Background()
	start := time.Now()
	ctx = server.OnBeforeStart(ctx, start)
	ctx = server.OnBeforeEnd(ctx, []attribute.KeyValue{}, start)
	server.OnAfterStart(ctx, start)
	server.OnAfterEnd(ctx, []attribute.KeyValue{}, time.Now())
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.server.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

func TestGlobalRpcClientMetrics(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	m := mp.Meter("test-meter")
	InitRpcMetrics(m)
	client := RpcClientMetrics("net.rpc.client")
	ctx := context.Background()
	start := time.Now()
	ctx = client.OnBeforeStart(ctx, start)
	ctx = client.OnBeforeEnd(ctx, []attribute.KeyValue{}, start)
	client.OnAfterStart(ctx, start)
	client.OnAfterEnd(ctx, []attribute.KeyValue{}, time.Now())
	rm := &metricdata.ResourceMetrics{}
	reader.Collect(ctx, rm)
	if rm.ScopeMetrics[0].Metrics[0].Name != "rpc.client.duration" {
		panic("wrong metrics name, " + rm.ScopeMetrics[0].Metrics[0].Name)
	}
}

func TestNilClientMeter(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	_ = metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	_, err := newRpcClientMetric("test", nil)
	if err == nil {
		panic(err)
	}
}

func TestNilServerMeter(t *testing.T) {
	reader := metric.NewManualReader()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	_ = metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))
	_, err := newRpcServerMetric("test", nil)
	if err == nil {
		panic(err)
	}
}
