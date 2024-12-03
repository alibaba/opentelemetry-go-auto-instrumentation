// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import "net/http"

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("echo"))
		if err != nil {
			panic(err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
