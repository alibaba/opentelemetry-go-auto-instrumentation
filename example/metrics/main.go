// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"
	"runtime"
)

func main() {
	http.Handle("/gc-metrics", http.HandlerFunc(gcMetrics))
	http.Handle("/mem-metrics", http.HandlerFunc(memMetrics))
	err := http.ListenAndServe("0.0.0.0:9000", nil)
	if err != nil {
		panic(err)
	}
}

func gcMetrics(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
	w.Write([]byte("Get GC Metrics"))
}

func memMetrics(w http.ResponseWriter, r *http.Request) {
	var bytes []byte
	for i := 0; i < 10; i++ {
		bytes = append(bytes, make([]byte, 1024*1024)...)
	}
	w.Write([]byte("Get Memory Metrics"))
}
