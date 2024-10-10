// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
