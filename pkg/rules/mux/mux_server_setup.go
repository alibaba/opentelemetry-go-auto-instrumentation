// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore

package rule

import (
	"bufio"
	"context"
	"fmt"
	mux "github.com/gorilla/mux"
	"net"
	"net/http"
	"strconv"
)

var muxInstrumenter = BuildMuxHttpServerOtelInstrumenter()

func muxServerOnEnter(call mux.CallContext, router *mux.Router, w http.ResponseWriter, req *http.Request) {
	muxRequest := muxHttpRequest{
		method:  req.Method,
		url:     req.URL,
		header:  req.Header,
		version: strconv.Itoa(req.ProtoMajor) + "." + strconv.Itoa(req.ProtoMinor),
		host:    req.Host,
		isTls:   req.TLS != nil,
	}
	ctx := muxInstrumenter.Start(req.Context(), muxRequest)
	x := call.GetParam(1).(http.ResponseWriter)
	x1 := &muxWriterWrapper{ResponseWriter: x, statusCode: http.StatusOK}
	call.SetParam(1, x1)
	call.SetParam(2, req.WithContext(ctx))
	call.SetKeyData("ctx", ctx)
	call.SetKeyData("request", muxRequest)
	return
}

func muxServerOnExit(call mux.CallContext) {
	c := call.GetKeyData("ctx")
	if c == nil {
		return
	}
	ctx, ok := c.(context.Context)
	if !ok {
		return
	}
	m := call.GetKeyData("request")
	if m == nil {
		return
	}
	muxRequest, ok := m.(muxHttpRequest)
	if !ok {
		return
	}
	if p, ok := call.GetParam(1).(http.ResponseWriter); ok {
		if w1, ok := p.(*muxWriterWrapper); ok {
			muxInstrumenter.End(ctx, muxRequest, muxHttpResponse{
				statusCode: w1.statusCode,
			}, nil)
		}
	}
	return
}

type muxWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *muxWriterWrapper) WriteHeader(statusCode int) {
	// cache the status code
	w.statusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *muxWriterWrapper) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("responseWriter does not implement http.Hijacker")
}
