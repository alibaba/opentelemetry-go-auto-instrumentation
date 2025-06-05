package dubbo

import (
	"context"
	"sync"

	_ "unsafe"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var dubboServerInstrumenter = BuildDubboServerInstrumenter()

var (
	dsf     *dubboServerOTelFilter
	dsfOnce sync.Once
)

func init() {
	extension.SetFilter(DubboServerOTelFilterKey, func() filter.Filter {
		if dsf == nil {
			dsfOnce.Do(func() {
				dsf = &dubboServerOTelFilter{
					Propagators: otel.GetTextMapPropagator(),
				}
			})
		}
		return dsf
	})
}

//go:linkname dubboNewServerOnEnter dubbo.apache.org/dubbo-go/v3.dubboNewServerOnEnter
func dubboNewServerOnEnter(call api.CallContext, instance *dubbo.Instance, opts ...server.ServerOption) {
	if !dubboEnabler.Enable() {
		return
	}
	opts = append(opts, server.WithServerFilter(DubboServerOTelFilterKey))
	call.SetParam(1, opts)
}

type dubboServerOTelFilter struct {
	Propagators propagation.TextMapPropagator
}

func (f *dubboServerOTelFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	if !dubboEnabler.Enable() {
		return invoker.Invoke(ctx, invocation)
	}

	attachments := invocation.Attachments()
	bags, spanCtx := extract(ctx, attachments, f.Propagators)
	ctx = baggage.ContextWithBaggage(ctx, bags)

	req := dubboRequest{
		methodName:    invocation.MethodName(),
		serviceKey:    invoker.GetURL().ServiceKey(),
		serverAddress: invoker.GetURL().Address(),
		attachments:   attachments,
	}

	ctx = dubboServerInstrumenter.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), req)

	result := invoker.Invoke(ctx, invocation)

	resp := dubboResponse{}
	if result.Error() != nil {
		resp.hasError = true
		resp.errorMsg = result.Error().Error()
	}

	dubboServerInstrumenter.End(ctx, req, resp, result.Error())

	return result
}

func (f *dubboServerOTelFilter) OnResponse(ctx context.Context, res protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	return res
}
