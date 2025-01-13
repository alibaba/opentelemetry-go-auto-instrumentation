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
	sendStreamReq(context.Background())
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.Assert(len(stubs) == 0, "Except client db system to be 0, got %d", len(stubs))
	}, 1)
}
