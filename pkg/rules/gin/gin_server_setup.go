//go:build ignore

package rule

import (
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

func nextOnEnter(call gin.CallContext, c *gin.Context) {
	if c == nil {
		return
	}
	val := reflect.ValueOf(*c)
	var handleName string
	index := val.FieldByName("index").Int() + 1
	if c.HandlerNames() != nil {
		handleName = c.HandlerNames()[index]
	} else {
		handleName = "gin-server"
	}
	ctx := netGinServerInstrument.Start(c.Request.Context(), ginRequest{
		method:     c.Request.Method,
		url:        *c.Request.URL,
		header:     c.Request.Header,
		handleName: handleName,
		version:    strconv.Itoa(c.Request.ProtoMajor) + "." + strconv.Itoa(c.Request.ProtoMinor),
		host:       c.Request.Host,
		isTls:      c.Request.TLS != nil,
	})
	data := make(map[string]interface{}, 1)
	data["ctx"] = ctx
	call.SetData(data)
	return
}

func nextOnExit(call gin.CallContext) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	if c, ok := call.GetParam(0).(*gin.Context); ok {
		statusCode := c.Writer.Status()
		netGinServerInstrument.End(ctx, ginRequest{}, ginResponse{
			statusCode: statusCode,
		}, nil)
	}

	return
}
