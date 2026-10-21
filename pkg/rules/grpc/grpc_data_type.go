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

package grpc

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	grpcTraceExporterPath = "/opentelemetry.proto.collector.trace.v1.TraceService/Export"
	grpcMetricExporterPath = "/opentelemetry.proto.collector.metrics.v1.MetricsService/Export"
)

var grpcClientInstrument = BuildGrpcClientInstrumenter()

type grpcRequest struct {
	methodName    string
	serverAddress string
	propagators   propagation.TextMapCarrier
}

type grpcResponse struct {
	statusCode int
}

type gRPCContextKey struct{}

type gRPCContext struct {
	methodName string
	// span is the span started by this handler in TagRPC. It is kept here
	// because, when other stats handlers (e.g. otelgrpc) are registered on the
	// same ClientConn/Server, trace.SpanFromContext(ctx) in HandleRPC returns
	// the span of the handler tagged last instead of ours.
	span trace.Span
}
