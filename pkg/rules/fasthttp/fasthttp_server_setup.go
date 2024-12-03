// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fasthttp

import (
	"net/url"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/valyala/fasthttp"
)

var fastHttpServerInstrumenter = BuildFastHttpServerOtelInstrumenter()

func newFastHttpServerDelegateHandler(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()
		handler(ctx)
		u, err := url.Parse(ctx.URI().String())
		if err != nil {
			return
		}
		request := fastHttpRequest{
			method: string(ctx.Method()),
			url:    u,
			isTls:  ctx.IsTLS(),
			header: &ctx.Request.Header,
		}
		fastHttpServerInstrumenter.StartAndEnd(ctx, request, fastHttpResponse{
			statusCode: ctx.Response.StatusCode(),
			header:     &ctx.Response.Header,
		}, ctx.Err(), startTime, time.Now())
	}
}

func listenAndServeFastHttpOnEnter(call api.CallContext, s *fasthttp.Server, addr string) {
	if !fastHttpEnabler.Enable() {
		return
	}
	if s == nil {
		return
	}
	s.Handler = newFastHttpServerDelegateHandler(s.Handler)
	call.SetParam(0, s)
}
