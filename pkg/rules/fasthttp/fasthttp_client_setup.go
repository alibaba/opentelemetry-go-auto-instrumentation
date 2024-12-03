// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fasthttp

import (
	"context"
	"net/url"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/valyala/fasthttp"
)

var fastHttpClientInstrumenter = BuildFastHttpClientOtelInstrumenter()

func clientFastHttpOnEnter(call api.CallContext, c *fasthttp.HostClient, req *fasthttp.Request, resp *fasthttp.Response) {
	if !fastHttpEnabler.Enable() {
		return
	}
	scheme := req.URI().Scheme()
	isTLS := false
	if string(scheme) == "https" {
		isTLS = true
	}
	u, err := url.Parse(req.URI().String())
	if err != nil {
		return
	}
	request := fastHttpRequest{
		method: string(req.Header.Method()),
		url:    u,
		isTls:  isTLS,
		header: &req.Header,
	}
	ctx := fastHttpClientInstrumenter.Start(context.Background(), request)
	data := make(map[string]interface{}, 3)
	data["ctx"] = ctx
	data["request"] = request
	data["response"] = resp
	call.SetData(data)
}

func clientFastHttpOnExit(call api.CallContext, err error) {
	if !fastHttpEnabler.Enable() {
		return
	}
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)
	request := data["request"].(fastHttpRequest)
	resp := data["response"].(*fasthttp.Response)
	fastHttpClientInstrumenter.End(ctx, request, fastHttpResponse{
		statusCode: resp.StatusCode(),
		header:     &resp.Header,
	}, err)
}
