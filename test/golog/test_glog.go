// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger.Info("slog logger")
	log.Printf("go log")
	w.Write([]byte("hello world"))
}

func main() {
	http.HandleFunc("/hello", hello)
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	time.Sleep(5 * time.Second)
	client := http.Client{}
	client.Get("http://localhost:8080/hello")
}
