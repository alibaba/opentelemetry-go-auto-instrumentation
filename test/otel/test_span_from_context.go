// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

func main() {
	go func() {
		http.HandleFunc("/otel", func(writer http.ResponseWriter, request *http.Request) {
			span := trace.SpanFromContext(context.Background())
			if !span.IsRecording() {
				panic("span should be recordedc")
			}
			if !span.SpanContext().IsValid() {
				panic("span should be valid")
			}
			fmt.Printf("%v\n", span)
			writer.Write([]byte("hello otel"))
		})
		http.ListenAndServe(":8989", nil)
	}()
	time.Sleep(3 * time.Second)
	client := http.Client{}
	client.Get("http://127.0.0.1:8989/otel")
}
