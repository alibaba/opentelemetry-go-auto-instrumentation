package main

import (
	"bytes"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"net/http"
	"strconv"
	"time"
)

func setupBasicHttp() {
	http.HandleFunc("/a", redirectHandler)
	http.HandleFunc("/b", helloHandler)
	var err error
	port, err = verifier.GetFreePort()
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	go setupBasicHttp()
	time.Sleep(1 * time.Second)
	_, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/a")
	if err != nil {
		panic(err)
	}
	jsonData := []byte(`{"key1": "value1", "key2": "value2"}`)
	req, err := http.NewRequest("POST", "http://127.0.0.1:"+strconv.Itoa(port)+"/a", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	time.Sleep(1 * time.Second)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:"+strconv.Itoa(port)+"/a", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), 200, 0, int64(port))
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /a", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), "Go-http-client/1.1", "", "/a", "", "/a", 200)
		verifier.VerifyHttpClientAttributes(stubs[0][2], "GET", "GET", "http://127.0.0.1:"+strconv.Itoa(port)+"/b", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), 200, 0, int64(port))
		verifier.VerifyHttpServerAttributes(stubs[0][3], "GET /b", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), "Go-http-client/1.1", "", "/b", "", "/b", 200)
		if stubs[0][1].Parent.TraceID().String() != stubs[0][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
		if stubs[0][2].Parent.TraceID().String() != stubs[0][1].SpanContext.TraceID().String() {
			log.Fatal("span 2 should be child of span 1")
		}
		if stubs[0][3].Parent.TraceID().String() != stubs[0][2].SpanContext.TraceID().String() {
			log.Fatal("span 3 should be child of span 2")
		}

		verifier.VerifyHttpClientAttributes(stubs[1][0], "POST", "POST", "http://127.0.0.1:"+strconv.Itoa(port)+"/a", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), 200, 0, int64(port))
		verifier.VerifyHttpServerAttributes(stubs[1][1], "POST /a", "POST", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), "Go-http-client/1.1", "", "/a", "", "/a", 200)
		verifier.VerifyHttpClientAttributes(stubs[1][2], "GET", "GET", "http://127.0.0.1:"+strconv.Itoa(port)+"/b", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), 200, 0, int64(port))
		verifier.VerifyHttpServerAttributes(stubs[1][3], "GET /b", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), "Go-http-client/1.1", "", "/b", "", "/b", 200)
		if stubs[1][1].Parent.TraceID().String() != stubs[1][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
		if stubs[1][2].Parent.TraceID().String() != stubs[1][1].SpanContext.TraceID().String() {
			log.Fatal("span 2 should be child of span 1")
		}
		if stubs[1][3].Parent.TraceID().String() != stubs[1][2].SpanContext.TraceID().String() {
			log.Fatal("span 3 should be child of span 2")
		}
	})
}
