//go:build ignore

package rule

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

func htmlOnEnter(call gin.CallContext, c *gin.Context, code int, name string, obj any) {
	if c == nil {
		return
	}
	ctx := netGinServerInstrument.Start(c.Request.Context(), ginRequest{
		method:  c.Request.Method,
		url:     *c.Request.URL,
		header:  c.Request.Header,
		version: strconv.Itoa(c.Request.ProtoMajor) + "." + strconv.Itoa(c.Request.ProtoMinor),
		host:    c.Request.Host,
		isTls:   c.Request.TLS != nil,
	})
	data := make(map[string]interface{}, 2)
	data["ctx"] = ctx
	data["code"] = code
	call.SetData(data)
	return
}

func htmlOnExit(call gin.CallContext) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	code := data["code"].(int)
	netGinServerInstrument.End(ctx, ginRequest{}, ginResponse{
		statusCode: code,
	}, nil)

	return
}
