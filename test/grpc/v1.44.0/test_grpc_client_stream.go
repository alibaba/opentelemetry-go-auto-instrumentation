package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"time"
)

func main() {
	// starter server
	go setupGRPC()
	time.Sleep(3 * time.Second)
	// use a http client to request to the server
	sendStreamReq(context.Background())
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		// filter out client
		verifier.Assert(len(stubs[0]) == 1, "Except client grpc system to be 1, got %d", len(stubs))
		verifier.Assert(stubs[0][0].SpanKind == trace.SpanKindServer, "Except client grpc system to be server, got %v", stubs[0][0].SpanKind)
	}, 1)
}
