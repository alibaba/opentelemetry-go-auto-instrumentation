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

package gin

import (
	"context"
	"reflect"
	"strconv"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/gin-gonic/gin"
)

func nextOnEnter(call api.CallContext, c *gin.Context) {
	if !ginEnabler.Enable() {
		return
	}
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

func nextOnExit(call api.CallContext) {
	if !ginEnabler.Enable() {
		return
	}
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
