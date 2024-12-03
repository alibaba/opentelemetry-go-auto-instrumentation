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

func countries(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("test"))
}

func setupPattern() {
	r := mux.NewRouter()
	r.HandleFunc("/{name}/countries/{country}", countries).Methods(http.MethodGet)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func main() {
	go setupPattern()
	time.Sleep(5 * time.Second)
	client := &http.Client{}
	resp, err := client.Get("http://127.0.0.1:8080/1/countries/2")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:8080/1/countries/2", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:8080", 200, 0, 8080)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "/{name}/countries/{country}", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8080", "Go-http-client/1.1", "http", "/1/countries/2", "", "/{name}/countries/{country}", 200)
	}, 1)
}
