package main

import (
	"context"
	"net/http"
)

var client *http.Client
var req *http.Request
var req1 *http.Request

func init() {
	client = &http.Client{}
	// create request ahead of time, to test that when the instrumentation still work
	req, _ = http.NewRequest("GET", "http://localhost:8080", nil)
	ctx := context.Background()
	req1, _ = http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
}

func main() {
	client.Do(req)
	client.Do(req1)
	e := &http.MaxBytesError{Limit: 0}
	msg := e.Error()
	println(e.Limit)
	println(msg)
}
