package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"time"
)

func setupHttp() {
	SetupGRPC()
}

func requestServer() {
	SendReq(context.Background())
}

func main() {
	// starter server
	go setupHttp()
	time.Sleep(3 * time.Second)
	// use a http client to request to the server
	requestServer()
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		x, _ := json.Marshal(stubs)
		fmt.Println(string(x))
		verifier.VerifyGrpcClientAttributes(stubs[0][0], "/HelloGrpc/Hello", "/HelloGrpc/Hello", "grpc", "tcp", "ipv4", "", "[::1]:9003", 200)
		verifier.VerifyGrpcServerAttributes(stubs[0][1], "/HelloGrpc/Hello", "/HelloGrpc/Hello", "grpc", "tcp", "ipv4", "", "", 200)
		if stubs[0][1].Parent.TraceID().String() != stubs[0][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
	}, 1)

}
