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

package rpc

import "testing"

type testRequest struct {
	System  string
	Service string
	Method  string
}

type testGetter struct {
}

func (t testGetter) GetSystem(request testRequest) string {
	if request.System != "" {
		return request.System
	}
	return ""
}

func (t testGetter) GetService(request testRequest) string {
	if request.Service != "" {
		return request.Service
	}
	return ""
}

func (t testGetter) GetMethod(request testRequest) string {
	if request.Method != "" {
		return request.Method
	}
	return ""
}

func TestExtractSpanName(t *testing.T) {
	r := RpcSpanNameExtractor[testRequest]{getter: testGetter{}}
	spanName := r.Extract(testRequest{Method: "method", Service: "service"})
	if spanName != "service/method" {
		t.Fatalf("extract span name extractor failed, expected 'service/method', got '%s'", spanName)
	}
	spanName = r.Extract(testRequest{})
	if spanName != "RPC request" {
		t.Fatalf("extract span name extractor failed, expected 'RPC request', got '%s'", spanName)
	}
}
