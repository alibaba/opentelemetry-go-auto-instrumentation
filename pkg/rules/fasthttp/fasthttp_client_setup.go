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
//go:build ignore

package fasthttp

import (
	"context"
	"github.com/valyala/fasthttp"
	"net/url"
)

var fastHttpClientInstrumenter = BuildFastHttpClientOtelInstrumenter()

func clientFastHttpOnEnter(call fasthttp.CallContext, c *fasthttp.HostClient, req *fasthttp.Request, resp *fasthttp.Response) {
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

func clientFastHttpOnExit(call fasthttp.CallContext, err error) {
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)
	request := data["request"].(fastHttpRequest)
	resp := data["response"].(*fasthttp.Response)
	fastHttpClientInstrumenter.End(ctx, request, fastHttpResponse{
		statusCode: resp.StatusCode(),
		header:     &resp.Header,
	}, err)
}
