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
	r := HttpClientSpanNameExtractor[testRequest, testResponse]{getter: testClientGetter{}}
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
	r := HttpServerSpanNameExtractor[testRequest, testResponse]{getter: testServerGetter{}}
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
