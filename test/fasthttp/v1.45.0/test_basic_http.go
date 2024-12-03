// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	go func() {
		fasthttp.ListenAndServe(":8080", hello)
	}()
	time.Sleep(5 * time.Second)
	client := &fasthttp.Client{}
	reqURL := "http://localhost:8080"
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.SetRequestURI(reqURL)
	req.Header.SetMethod(fasthttp.MethodGet)
	err := client.Do(req, resp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://localhost:8080/", "http", "", "tcp", "ipv4", "", "localhost:8080", 200, 0, 8080)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /", "GET", "http", "tcp", "ipv4", "", "localhost:8080", "fasthttp", "http", "/", "", "/", 200)
	}, 1)
}
