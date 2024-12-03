// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"math"
	"net"
	"net/http"
	"os"
	"time"
)

var providerIp = os.Getenv("OTEL_PROVIDER_IP")

var netTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 2 * time.Minute,
	}).DialContext,
	MaxIdleConns:          0,
	MaxIdleConnsPerHost:   math.MaxInt,
	IdleConnTimeout:       90 * time.Second,
	ExpectContinueTimeout: 10 * time.Second,
}

var customClient = http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

func setup() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		resp, err := customClient.Get("http://" + providerIp + ":8080/echo")
		defer func() {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}()
		if err != nil {
			panic(err)
		}
		_, err = w.Write([]byte("echo"))
		if err != nil {
			panic(err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func main() {
	setup()
}
