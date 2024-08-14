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
