//go:build ignore

package rule

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

var netHttpServerInstrumenter = BuildNetHttpServerOtelInstrumenter()

func serverOnEnter(call http.CallContext, _ interface{}, w http.ResponseWriter, r *http.Request) {
	request := netHttpRequest{
		method:  r.Method,
		url:     *r.URL,
		header:  r.Header,
		version: strconv.Itoa(r.ProtoMajor) + "." + strconv.Itoa(r.ProtoMinor),
		host:    r.Host,
		isTls:   r.TLS != nil,
	}
	ctx := netHttpServerInstrumenter.Start(r.Context(), request)
	if x, ok := call.GetParam(1).(http.ResponseWriter); ok {
		x1 := &writerWrapper{ResponseWriter: x, statusCode: http.StatusOK}
		call.SetParam(1, x1)
	}
	call.SetParam(2, r.WithContext(ctx))
	data := make(map[string]interface{}, 1)
	data["ctx"] = ctx
	data["request"] = request
	call.SetData(data)
	return
}

func serverOnExit(call http.CallContext) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	request, ok := data["request"].(netHttpRequest)
	if !ok {
		return
	}
	if p, ok := call.GetParam(1).(http.ResponseWriter); ok {
		if w1, ok := p.(*writerWrapper); ok {
			netHttpServerInstrumenter.End(ctx, request, netHttpResponse{
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
	w.statusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *writerWrapper) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("responseWriter does not implement http.Hijacker")
}
