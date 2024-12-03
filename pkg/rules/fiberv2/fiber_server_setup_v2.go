// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fiberv2

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/url"
)

var fiberv2ServerInstrumenter = BuildFiberV2ServerOtelInstrumenter()

func fiberHttpOnEnterv2(call api.CallContext, app *fiber.App, ctx *fasthttp.RequestCtx) {
	if !fiberv2Enabler.Enable() {
		return
	}
	u, err := url.Parse(ctx.URI().String())
	if err != nil {
		return
	}
	request := &fiberv2Request{
		method: string(ctx.Method()),
		url:    u,
		isTls:  ctx.IsTLS(),
		header: &ctx.Request.Header,
	}
	ctxSpan := fiberv2ServerInstrumenter.Start(ctx, request)
	data := make(map[string]interface{}, 2)
	data["ctx"] = ctx
	data["ctxSpan"] = ctxSpan
	data["request"] = request
	call.SetData(data)
	return
}

func fiberHttpOnExitv2(call api.CallContext) {
	if call.GetData() == nil {
		return
	}
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(*fasthttp.RequestCtx)
	ctxSpan := data["ctxSpan"].(context.Context)
	request, ok := data["request"].(*fiberv2Request)
	if !ok {
		return
	}
	fiberv2ServerInstrumenter.End(ctxSpan, request, &fiberv2Response{
		statusCode: ctx.Response.StatusCode(),
		header:     &ctx.Response.Header,
	}, nil)

}
