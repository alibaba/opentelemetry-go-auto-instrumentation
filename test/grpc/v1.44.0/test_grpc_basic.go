// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"time"
)

func main() {
	// starter server
	go setupGRPC()
	time.Sleep(3 * time.Second)
	// use a http client to request to the server
	sendReq(context.Background())
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyRpcClientAttributes(stubs[0][0], "/HelloGrpc/Hello", "grpc", "/HelloGrpc", "Hello")
		verifier.VerifyRpcServerAttributes(stubs[0][1], "/HelloGrpc/Hello", "grpc", "/HelloGrpc", "Hello")
	}, 1)
}
