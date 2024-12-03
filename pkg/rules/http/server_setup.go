// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

var netHttpServerInstrumenter = BuildNetHttpServerOtelInstrumenter()

func serverOnEnter(call api.CallContext, _ interface{}, w http.ResponseWriter, r *http.Request) {
	if netHttpFilter.FilterUrl(r.URL) {
		return
	}
	request := &netHttpRequest{
		method:  r.Method,
		url:     r.URL,
		header:  r.Header,
		version: getProtocolVersion(r.ProtoMajor, r.ProtoMinor),
		host:    r.Host,
		isTls:   r.TLS != nil,
	}
	ctx := netHttpServerInstrumenter.Start(r.Context(), request)
	if x, ok := call.GetParam(1).(http.ResponseWriter); ok {
		x1 := &writerWrapper{ResponseWriter: x, statusCode: http.StatusOK}
		call.SetParam(1, x1)
	}
	call.SetParam(2, r.WithContext(ctx))
	data := make(map[string]interface{}, 2)
	data["ctx"] = ctx
	data["request"] = request
	call.SetData(data)
	return
}

func serverOnExit(call api.CallContext) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	request, ok := data["request"].(*netHttpRequest)
	if !ok {
		return
	}
	if p, ok := call.GetParam(1).(http.ResponseWriter); ok {
		if w1, ok := p.(*writerWrapper); ok {
			netHttpServerInstrumenter.End(ctx, request, &netHttpResponse{
				statusCode: w1.statusCode,
			}, nil)
		}
	}

	return
}

type writerWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *writerWrapper) WriteHeader(statusCode int) {
	// cache the status code
	if w.statusCode == statusCode {
		return // 防止多次写入 Header
	}
	w.statusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *writerWrapper) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("responseWriter does not implement http.Hijacker")
}

func (w *writerWrapper) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *writerWrapper) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (w *writerWrapper) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
