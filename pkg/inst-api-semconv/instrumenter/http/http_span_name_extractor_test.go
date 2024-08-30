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

package http

import "testing"

type testRequest struct {
	Method string
	Route  string
}

type testResponse struct {
}

type testClientGetter struct {
	HttpClientAttrsGetter[testRequest, testResponse]
}

type testServerGetter struct {
	HttpServerAttrsGetter[testRequest, testResponse]
}

func (t testClientGetter) GetRequestMethod(request testRequest) string {
	if request.Method != "" {
		return request.Method
	}
	return ""
}

func (t testClientGetter) GetSpanName(request testRequest) string {
	if request.Method != "" {
		return request.Method
	}
	return "HTTP"
}

func (t testServerGetter) GetRequestMethod(request testRequest) string {
	if request.Method != "" {
		return request.Method
	}
	return ""
}

func (t testServerGetter) GetHttpRoute(request testRequest) string {
	if request.Route != "" {
		return request.Route
	}
	return ""
}

func TestHttpClientExtractSpanName(t *testing.T) {
	r := HttpClientSpanNameExtractor[testRequest, testResponse]{Getter: testClientGetter{}}
	spanName := r.Extract(testRequest{Method: "GET"})
	if spanName != "GET" {
		t.Errorf("want GET, got %s", spanName)
	}
	spanName = r.Extract(testRequest{})
	if spanName != "HTTP" {
		t.Errorf("want HTTP, got %s", spanName)
	}
}

func TestHttpServerExtractSpanName(t *testing.T) {
	r := HttpServerSpanNameExtractor[testRequest, testResponse]{Getter: testServerGetter{}}
	spanName := r.Extract(testRequest{Method: "GET"})
	if spanName != "GET" {
		t.Errorf("want GET, got %s", spanName)
	}
	spanName = r.Extract(testRequest{})
	if spanName != "HTTP" {
		t.Errorf("want HTTP, got %s", spanName)
	}
	spanName = r.Extract(testRequest{Method: "GET", Route: "/a/b"})
	if spanName != "GET /a/b" {
		t.Errorf("want GET /a/b, got %s", spanName)
	}
}
