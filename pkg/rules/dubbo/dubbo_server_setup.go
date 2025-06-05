package dubbo

import (
	"context"
	_ "unsafe"

	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

var dubboServerInstrumenter = BuildDubboServerInstrumenter()

//go:linkname dubboProviderGracefulShutdownFilterInvokeOnEnter dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown.dubboProviderGracefulShutdownFilterInvokeOnEnter
func dubboProviderGracefulShutdownFilterInvokeOnEnter(call api.CallContext, _ interface{}, ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) {
	if !dubboEnabler.Enable() {
		return
	}

	attachments := invocation.Attachments()
	bags, spanCtx := extract(ctx, attachments, otel.GetTextMapPropagator())
	ctx = baggage.ContextWithBaggage(ctx, bags)

	req := dubboRequest{
		methodName:    invocation.MethodName(),
		serviceKey:    invoker.GetURL().ServiceKey(),
		serverAddress: invoker.GetURL().Address(),
		attachments:   attachments,
	}

	ctx = dubboServerInstrumenter.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), req)

	data := make(map[string]interface{})
	data["ctx"] = ctx
	data["req"] = req

	call.SetData(data)
}

//go:linkname dubboProviderGracefulShutdownFilterInvokeOnExit dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown.dubboProviderGracefulShutdownFilterInvokeOnExit
func dubboProviderGracefulShutdownFilterInvokeOnExit(call api.CallContext, res protocol.Result) {
	if !dubboEnabler.Enable() {
		return
	}

	data := call.GetData().(map[string]interface{})
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	req, ok := data["req"].(dubboRequest)
	if !ok {
		return
	}

	resp := dubboResponse{}
	if res.Error() != nil {
		resp.hasError = true
		resp.errorMsg = res.Error().Error()
	}

	dubboServerInstrumenter.End(ctx, req, resp, res.Error())
}
