// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	st "dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

var dubboServerInstrument = BuildDubboServerInstrumenter()

func DubboServerOnEnter(call st.CallContext, _ interface{}, ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) {
	url := invoker.GetURL()
	if url == nil {
		return
	}
	attachments := invocation.Attachments()
	if attachments == nil {
		attachments = map[string]interface{}{}
	}
	req := dubboRequest{
		method:   GenerateSpanName(invoker, invocation),
		addr:     "",
		metadata: attachments,
	}
	ctxD := dubboServerInstrument.Start(ctx, req)
	data := make(map[string]interface{}, 3)
	data["ctx"] = ctxD
	data["dubboRequest"] = req
	call.SetData(data)
	return
}

func DubboServerOnExit(call st.CallContext, r protocol.Result) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	req := data["dubboRequest"].(dubboRequest)
	if r.Error() != nil {
		dubboServerInstrument.End(ctx, req, dubboResponse{
			statusCode: "500",
		}, r.Error())
	} else {
		dubboServerInstrument.End(ctx, req, dubboResponse{
			statusCode: "200",
		}, nil)
	}
}
