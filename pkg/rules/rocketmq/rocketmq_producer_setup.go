package rocketmq

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"reflect"
	"unsafe"
)

// SendSync埋点入口函数
func producerSendSyncOnEnter(call api.CallContext, _ interface{}, ctx context.Context, msg *primitive.Message, res *primitive.SendResult) {
	if !rocketmqEnabler.Enable() {
		return
	}
	//client := reflect.ValueOf(producer).Elem().FieldByName("client")
	req := rocketmqProducerReq{
		message: msg,
	}
	newCtx := producerInstrumenter.Start(ctx, req)
	data := make(map[string]interface{}, 3)
	data["ctx"] = newCtx
	data["res"] = res
	call.SetData(data)
}

// SendSync埋点出口函数
func producerSendSyncOnExit(call api.CallContext, err error) {
	if !rocketmqEnabler.Enable() {
		return
	}

	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)

	// 构建请求对象
	req := rocketmqProducerReq{
		message: call.GetParam(2).(*primitive.Message),
	}

	// 提前检查响应是否为空
	resObj, ok := data["res"]
	if !ok || resObj == nil {
		producerInstrumenter.End(ctx, req, rocketmqProducerRes{}, err)
		return
	}

	res, ok := resObj.(*primitive.SendResult)
	if !ok || res == nil || res.MessageQueue == nil || err != nil {
		producerInstrumenter.End(ctx, req, rocketmqProducerRes{result: res}, err)
		return
	}

	// 获取客户端对象
	clientValue := reflect.ValueOf(call.GetParam(0)).Elem().FieldByName("client")
	clientValue = reflect.NewAt(clientValue.Type(), unsafe.Pointer(clientValue.UnsafeAddr())).Elem()

	// 如果客户端对象无效，直接结束埋点
	if !clientValue.IsValid() || !clientValue.CanAddr() {
		producerInstrumenter.End(ctx, req, rocketmqProducerRes{result: res}, err)
		return
	}

	// 获取broker地址
	addr := getBrokerAddr(clientValue, res)

	// 获取客户端ID
	req.clientID = getClientID(clientValue)

	// 结束埋点
	producerInstrumenter.End(ctx, req, rocketmqProducerRes{result: res, addr: addr}, err)
}

// SendAsync埋点入口函数
func producerSendAsyncOnEnter(call api.CallContext, _ interface{}, ctx context.Context, msg *primitive.Message, originalCallback func(context.Context, *primitive.SendResult, error)) {
	if !rocketmqEnabler.Enable() {
		return
	}
	req := rocketmqProducerReq{
		message: msg,
	}

	// 启动埋点
	newCtx := producerInstrumenter.Start(ctx, req)
	// 包装回调函数
	wrappedCallback := func(callbackCtx context.Context, result *primitive.SendResult, callbackErr error) {
		// 调用原始回调
		originalCallback(callbackCtx, result, callbackErr)

		// 获取客户端对象
		clientValue := reflect.ValueOf(call.GetParam(0)).Elem().FieldByName("client")
		clientValue = reflect.NewAt(clientValue.Type(), unsafe.Pointer(clientValue.UnsafeAddr())).Elem()

		// 如果客户端对象无效或结果为空，直接结束埋点
		if !clientValue.IsValid() || !clientValue.CanAddr() || result == nil || result.MessageQueue == nil {
			producerInstrumenter.End(ctx, req, rocketmqProducerRes{result: result}, callbackErr)
			return
		}

		// 获取broker地址
		addr := getBrokerAddr(clientValue, result)

		// 获取客户端ID
		req.clientID = getClientID(clientValue)
		res := rocketmqProducerRes{
			result: result,
			addr:   addr,
		}
		// 结束埋点
		producerInstrumenter.End(newCtx, req, res, callbackErr)
	}
	// 替换原始回调为包装后的回调
	call.SetParam(3, wrappedCallback)

}

// SendOneWay 埋点入口函数
func producerSendOneWayOnEnter(call api.CallContext, _ interface{}, ctx context.Context, msg *primitive.Message) {
	if !rocketmqEnabler.Enable() {
		return
	}

	// 创建请求信息
	req := rocketmqProducerReq{
		message: msg,
	}

	// 开始埋点
	newCtx := producerInstrumenter.Start(ctx, req)

	// 保存上下文和请求信息
	data := map[string]interface{}{
		"ctx": newCtx,
		"req": req,
	}
	call.SetData(data)
}

// SendOneWay 埋点出口函数
func producerSendOneWayOnExit(call api.CallContext, err error) {
	if !rocketmqEnabler.Enable() {
		return
	}
	// 获取入口函数保存的上下文和请求
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)
	req := data["req"].(rocketmqProducerReq)

	// 客户端值检查
	clientValue := reflect.ValueOf(call.GetParam(0)).Elem().FieldByName("client")
	clientValue = reflect.NewAt(clientValue.Type(), unsafe.Pointer(clientValue.UnsafeAddr())).Elem()

	// 如果客户端对象无效或消息队列为空，直接结束埋点
	if !clientValue.IsValid() || !clientValue.CanAddr() ||
		req.message == nil || req.message.Queue == nil {
		producerInstrumenter.End(ctx, req, rocketmqProducerRes{}, err)
		return
	}

	// 获取broker地址
	addr := getOneWayBrokerAddr(clientValue, req.message)

	// 获取客户端ID
	req.clientID = getClientID(clientValue)

	// 结束埋点
	producerInstrumenter.End(ctx, req, rocketmqProducerRes{addr: addr}, err)
}

// 从客户端获取broker地址
func getBrokerAddr(clientValue reflect.Value, res *primitive.SendResult) string {
	if res == nil || res.MessageQueue == nil || res.MessageQueue.BrokerName == "" {
		return ""
	}

	getNameSrvMethod := clientValue.MethodByName("GetNameSrv")
	if !getNameSrvMethod.IsValid() {
		return ""
	}

	nameSrvResults := getNameSrvMethod.Call(nil)
	if len(nameSrvResults) == 0 || !nameSrvResults[0].IsValid() {
		return ""
	}

	nameSrv := nameSrvResults[0]
	findMethod := nameSrv.MethodByName("FindBrokerAddrByName")
	if !findMethod.IsValid() {
		return ""
	}

	findResults := findMethod.Call([]reflect.Value{
		reflect.ValueOf(res.MessageQueue.BrokerName),
	})

	if len(findResults) == 0 || !findResults[0].IsValid() {
		return ""
	}

	return findResults[0].String()
}

// 获取客户端ID
func getClientID(clientValue reflect.Value) string {
	getClientIDMethod := clientValue.MethodByName("ClientID")
	if !getClientIDMethod.IsValid() {
		return ""
	}

	clientIDResults := getClientIDMethod.Call(nil)
	if len(clientIDResults) == 0 || !clientIDResults[0].IsValid() {
		return ""
	}

	return clientIDResults[0].String()
}

// 从客户端获取SendOneWay模式下的broker地址
func getOneWayBrokerAddr(clientValue reflect.Value, msg *primitive.Message) string {
	if msg == nil || msg.Queue == nil || msg.Queue.BrokerName == "" {
		return ""
	}

	getNameSrvMethod := clientValue.MethodByName("GetNameSrv")
	if !getNameSrvMethod.IsValid() {
		return ""
	}

	nameSrvResults := getNameSrvMethod.Call(nil)
	if len(nameSrvResults) == 0 || !nameSrvResults[0].IsValid() {
		return ""
	}

	nameSrv := nameSrvResults[0]
	findMethod := nameSrv.MethodByName("FindBrokerAddrByName")
	if !findMethod.IsValid() {
		return ""
	}

	findResults := findMethod.Call([]reflect.Value{
		reflect.ValueOf(msg.Queue.BrokerName),
	})

	if len(findResults) == 0 || !findResults[0].IsValid() {
		return ""
	}

	return findResults[0].String()
}
