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

package dubbo

import "testing"

type testRequest struct {
	Method string
	Addr   string
}

type testResponse struct {
	statusCode string
}

type testClientGetter struct {
	DubboClientAttrsGetter[testRequest, testResponse]
}

type testServerGetter struct {
	DubboServerAttrsGetter[testRequest, testResponse]
}

func (t testClientGetter) GetRequestMethod(request testRequest) string {
	if request.Method != "" {
		return request.Method
	}
	return ""
}

func (t testServerGetter) GetRequestMethod(request testRequest) string {
	if request.Method != "" {
		return request.Method
	}
	return ""
}

func TestDubboClientExtractSpanName(t *testing.T) {
	r := DubboClientSpanNameExtractor[testRequest, testResponse]{Getter: testClientGetter{}}
	spanName := r.Extract(testRequest{Method: "/org.apache.dubbogo.samples.api.Greeter/SayHello"})
	if spanName != "/org.apache.dubbogo.samples.api.Greeter/SayHello" {
		t.Errorf("want /org.apache.dubbogo.samples.api.Greeter/SayHello, got %s", spanName)
	}
	spanName = r.Extract(testRequest{})
	if spanName != "DUBBO" {
		t.Errorf("want DUBBO, got %s", spanName)
	}
}

func TestDubboServerExtractSpanName(t *testing.T) {
	r := DubboServerSpanNameExtractor[testRequest, testResponse]{Getter: testServerGetter{}}
	spanName := r.Extract(testRequest{Method: "GET"})
	if spanName != "GET" {
		t.Errorf("want GET, got %s", spanName)
	}
	spanName = r.Extract(testRequest{})
	if spanName != "DUBBO" {
		t.Errorf("want HTTP, got %s", spanName)
	}
}
