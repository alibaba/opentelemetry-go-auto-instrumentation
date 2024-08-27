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
