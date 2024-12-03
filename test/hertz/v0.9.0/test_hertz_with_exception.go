// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"time"
)

func main() {
	go setupWithException()
	time.Sleep(5 * time.Second)
	GetException()
	time.Sleep(1 * time.Second)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:8888/exception", "http", "", "tcp", "ipv4", "", "127.0.0.1:8888", 500, 0, 8888)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /exception", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8888", "Host", "http", "/exception", "", "/exception", 500)
		fmt.Printf("%v %v\n", stubs[0][0].Status.Code, stubs[0][1].Status.Code)
		if stubs[0][0].Status.Code != codes.Error {
			panic("span should be error state")
		}
	}, 1)
}
