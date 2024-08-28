package main

import (
	"encoding/json"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"strconv"
	"time"
)

func main() {
	// starter server
	go RunServer()
	time.Sleep(3 * time.Second)
	// use a http client to request to the server
	RunClient()
	time.Sleep(5 * time.Second)
	fmt.Println("http")
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		x, _ := json.Marshal(stubs)
		fmt.Println(string(x))
		var testTime int
		for _, stup := range stubs[0] {
			component := verifier.GetAttribute(stup.Attributes, "component.name").AsString()
			if component == "kratos-http-client" {
				verifier.VerifyKratosClientAttributes(stup, "/helloworld.v1.Greeter/SayHello", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(8777), 200, 0, int64(8777))
				testTime++
			}
			if component == "kratos-http-server" {
				verifier.VerifyKratosServerAttributes(stup, "/helloworld/client", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(8777), 200)
				testTime++
			}
		}
		if stubs[0][1].Parent.TraceID().String() != stubs[0][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
		if testTime != 2 {
			verifier.Assert(false, "test count should be 2")
		}
	}, 1)

	RunGrpcClient()
	time.Sleep(5 * time.Second)
	fmt.Println("grpc")
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		x, _ := json.Marshal(stubs)
		fmt.Println(string(x))
		var testTime int
		for _, stup := range stubs[0] {
			component := verifier.GetAttribute(stup.Attributes, "component.name").AsString()
			if component == "kratos-grpc-client" {
				verifier.VerifyKratosGrpcClientAttributes(stup, "/helloworld.v1.Greeter/SayHello", "grpc", "tcp", "ipv4", "", 200)
				testTime++
			}
			if component == "kratos-grpc-server" {
				verifier.VerifyKratosGrpcServerAttributes(stup, "/helloworld.v1.Greeter/SayHello", "grpc", "tcp", "ipv4", "", 200)
				testTime++
			}
		}
		if stubs[0][1].Parent.TraceID().String() != stubs[0][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
		if testTime != 2 {
			verifier.Assert(false, "test grpc count should be 2")
		}
	}, 1)

}
