// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"go.opentelemetry.io/otel/propagation"
)

var grpcClientInstrument = BuildGrpcClientInstrumenter()

type grpcRequest struct {
	methodName  string
	propagators propagation.TextMapCarrier
}

type grpcResponse struct {
	statusCode int
}

type gRPCContextKey struct{}

type gRPCContext struct {
	methodName string
}
