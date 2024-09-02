package main

import (
	"encoding/json"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	//"helloworld/server/client"
	"time"
)

import (
	//"dubbo.apache.org/dubbo-go/v3/common/logger"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

func main() {
	// starter server
	go RunServer()
	time.Sleep(3 * time.Second)
	RunClient()
	// use a http client to request to the server
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		x, _ := json.Marshal(stubs)
		fmt.Println(string(x))
		verifier.VerifyDubboClientAttributes(stubs[0][0], "/org.apache.dubbogo.samples.api.Greeter/SayHello", "/org.apache.dubbogo.samples.api.Greeter/SayHello", "dubbo", "tcp", "ipv4", "", "localhost:20000", 200)
		verifier.VerifyDubboServerAttributes(stubs[0][1], "/org.apache.dubbogo.samples.api.Greeter/SayHello", "/org.apache.dubbogo.samples.api.Greeter/SayHello", "dubbo", "tcp", "ipv4", "", "", 200)
		if stubs[0][1].Parent.TraceID().String() != stubs[0][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
	}, 1)
}
