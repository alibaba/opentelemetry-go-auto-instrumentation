package main

import (
	"context"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	// starter server
	go setupDubbo()
	time.Sleep(3 * time.Second)
	// use a client to request to the server
	sendErrDubboReq(context.Background())
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyRpcClientAttributes(stubs[0][0], "greet.GreetService/Greet", "apache_dubbo", "greet.GreetService", "Greet")
		verifier.VerifyRpcServerAttributes(stubs[0][1], "greet.GreetService/Greet", "apache_dubbo", "greet.GreetService", "Greet")
		if stubs[0][0].Status.Code != codes.Error {
			panic("wrong span status on span 0")
		}
		if stubs[0][1].Status.Code != codes.Error {
			panic("wrong span status on span 1")
		}
	}, 1)
}
