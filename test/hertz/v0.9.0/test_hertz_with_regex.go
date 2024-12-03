// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"time"
)

func main() {
	go setupWithRoute()
	time.Sleep(5 * time.Second)
	GetRoute()
	time.Sleep(1 * time.Second)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:8888/hertz/v1", "http", "", "tcp", "ipv4", "", "127.0.0.1:8888", 301, 0, 8888)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /hertz/v1/", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8888", "Host", "http", "/hertz/v1/", "", "/hertz/v1/", 301)
		verifier.VerifyHttpClientAttributes(stubs[1][0], "GET", "GET", "http://127.0.0.1:8888/hertz/v1/", "http", "", "tcp", "ipv4", "", "127.0.0.1:8888", 200, 0, 8888)
		verifier.VerifyHttpServerAttributes(stubs[1][1], "GET /hertz/v1/", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8888", "Host", "http", "/hertz/v1/", "", "/hertz/v1/", 200)
		verifier.VerifyHttpClientAttributes(stubs[2][0], "GET", "GET", "http://127.0.0.1:8888/hertz/v2/send", "http", "", "tcp", "ipv4", "", "127.0.0.1:8888", 200, 0, 8888)
		verifier.VerifyHttpServerAttributes(stubs[2][1], "GET /hertz/v2/send", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8888", "Host", "http", "/hertz/v2/send", "", "/hertz/v2/send", 200)
	}, 3)
}
