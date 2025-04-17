// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package rocketmq

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"reflect"
	"unsafe"
)

// consumeInner埋点入口函数
func consumerConsumeInnerOnEnter(call api.CallContext, _ interface{}, ctx context.Context, msgs []*primitive.MessageExt) {
	if !rocketmqEnabler.Enable() {
		return
	}

	// 检查消息数组是否为空
	if msgs == nil || len(msgs) == 0 {
		// 空消息情况下保存上下文
		newCtx := Start(ctx, msgs, "")
		data := map[string]interface{}{
			"ctx": newCtx,
			"req": msgs,
		}
		call.SetData(data)
		return
	}

	// 获取客户端对象
	clientValue := reflect.ValueOf(call.GetParam(0)).Elem().FieldByName("client")
	clientValue = reflect.NewAt(clientValue.Type(), unsafe.Pointer(clientValue.UnsafeAddr())).Elem()

	// 获取broker地址
	addr := getConsumerBrokerAddr(clientValue, msgs)

	// 启动埋点
	newCtx := Start(ctx, msgs, addr)

	// 保存上下文和请求信息
	data := map[string]interface{}{
		"ctx": newCtx,
		"req": msgs,
	}
	call.SetData(data)
}

// consumeInner埋点出口函数
func consumerConsumeInnerOnExit(call api.CallContext, consumeResult consumer.ConsumeResult, err error) {
	if !rocketmqEnabler.Enable() {
		return
	}

	// 安全获取入口函数保存的数据
	data, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}

	ctx, okCtx := data["ctx"].(context.Context)
	if !okCtx {
		return
	}

	req, okReq := data["req"].([]*primitive.MessageExt)
	if !okReq {
		End(ctx, nil, consumeResult, err)
		return
	}

	// 结束埋点
	End(ctx, req, consumeResult, err)
}

// 从客户端获取consumer的broker地址
func getConsumerBrokerAddr(clientValue reflect.Value, msgs []*primitive.MessageExt) string {
	// 检查客户端和消息的有效性
	if !clientValue.IsValid() || !clientValue.CanAddr() ||
		msgs == nil || len(msgs) == 0 ||
		msgs[0] == nil || msgs[0].Queue == nil ||
		msgs[0].Queue.BrokerName == "" {
		return ""
	}

	// 获取NameSrv
	getNameSrvMethod := clientValue.MethodByName("GetNameSrv")
	if !getNameSrvMethod.IsValid() {
		return ""
	}

	nameSrvResults := getNameSrvMethod.Call(nil)
	if len(nameSrvResults) == 0 || !nameSrvResults[0].IsValid() {
		return ""
	}

	// 获取Broker地址
	nameSrv := nameSrvResults[0]
	findMethod := nameSrv.MethodByName("FindBrokerAddrByName")
	if !findMethod.IsValid() {
		return ""
	}

	findResults := findMethod.Call([]reflect.Value{
		reflect.ValueOf(msgs[0].Queue.BrokerName),
	})

	if len(findResults) == 0 || !findResults[0].IsValid() {
		return ""
	}

	return findResults[0].String()
}
