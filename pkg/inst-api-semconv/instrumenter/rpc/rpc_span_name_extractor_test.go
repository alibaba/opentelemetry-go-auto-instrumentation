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
