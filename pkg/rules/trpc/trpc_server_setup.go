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

package trpc

import (
	"context"

	"trpc.group/trpc-go/trpc-go/codec"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

var trpcServerInstrumenter = BuildTrpcServerInstrumenter()

// func (s *service) handle(ctx context.Context, msg codec.Msg, reqBodyBuf []byte) (interface{}, error)
func serverTrpcOnEnter(call api.CallContext, _ interface{}, ctx context.Context, msg codec.Msg, reqBodyBuf []byte) {
	if !trpcEnabler.Enable() {
		return
	}
	request := trpcReq{
		callerMethod:  msg.CallerMethod(),
		callerService: msg.CallerService(),
		calleeService: msg.CalleeService(),
		calleeMethod:  msg.CalleeMethod(),
		msg:           msg,
	}
	newCtx := trpcServerInstrumenter.Start(ctx, request)
	data := make(map[string]interface{}, 3)
	data["ctx"] = newCtx
	data["request"] = request
	data["msg"] = msg
	call.SetData(data)
}

func serverTrpcOnExit(call api.CallContext, _ interface{}, err error) {
	if !trpcEnabler.Enable() {
		return
	}
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)
	request := data["request"].(trpcReq)
	msg := data["msg"].(codec.Msg)
	statusCode := 0
	if msg.ServerRspErr() != nil {
		statusCode = int(msg.ServerRspErr().Code)
	}
	trpcServerInstrumenter.End(ctx, request, trpcRes{
		stausCode: statusCode,
		msg:       msg,
	}, err)
}
