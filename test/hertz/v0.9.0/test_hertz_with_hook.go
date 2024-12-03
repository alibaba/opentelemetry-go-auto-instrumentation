// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"time"
)

func main() {
	go setupWithTracer()
	time.Sleep(5 * time.Second)
	GetDeadline()
	time.Sleep(1 * time.Second)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:8888/ping", "http", "", "tcp", "ipv4", "", "127.0.0.1:8888", 200, 0, 8888)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /ping", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8888", "Host", "http", "/ping", "", "/ping", 200)
	}, 1)
}
