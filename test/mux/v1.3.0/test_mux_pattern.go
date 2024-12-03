// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"io/ioutil"
	"net/http"
	"time"
)

func setupPattern() {
	r := mux.NewRouter()
	s := r.PathPrefix("/test").Subrouter()
	s.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func main() {
	go setupPattern()
	time.Sleep(5 * time.Second)
	client := &http.Client{}
	resp, err := client.Get("http://127.0.0.1:8080/test/1")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:8080/test/1", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:8080", 200, 0, 8080)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "/test/{key}", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8080", "Go-http-client/1.1", "http", "/test/1", "", "/test/{key}", 200)
	}, 1)
}
