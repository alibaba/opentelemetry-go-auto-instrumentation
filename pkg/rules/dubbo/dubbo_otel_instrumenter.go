package dubbo

import (
	"os"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/rpc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/trace"
)

type dubboEnable struct {
	enable bool
}

func (d *dubboEnable) Enable() bool {
	return d.enable
}

var dubboEnabler = dubboEnable{os.Getenv("OTEL_INSTRUMENTATION_DUBBO_ENABLED") != "false"}

type dubboAttrsGetter struct{}

func (g dubboAttrsGetter) GetSystem(request dubboRequest) string {
	return "apache_dubbo"
}

func (g dubboAttrsGetter) GetService(request dubboRequest) string {
	return request.serviceKey
}

func (g dubboAttrsGetter) GetMethod(request dubboRequest) string {
	return request.methodName
}

func (g dubboAttrsGetter) GetServerAddress(request dubboRequest) string {
	return request.serverAddress
}

type dubboStatusExtractor[REQUEST dubboRequest, RESPONSE dubboResponse] struct{}

func (g dubboStatusExtractor[REQUEST, RESPONSE]) Extract(span trace.Span, request dubboRequest, response dubboResponse, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else if response.hasError {
		span.SetStatus(codes.Error, response.errorMsg)
	} else {
		span.SetStatus(codes.Ok, codes.Ok.String())
	}
}

func BuildDubboClientInstrumenter() instrumenter.Instrumenter[dubboRequest, dubboResponse] {
	builder := instrumenter.Builder[dubboRequest, dubboResponse]{}
	clientGetter := dubboAttrsGetter{}
	return builder.Init().
		SetSpanStatusExtractor(&dubboStatusExtractor[dubboRequest, dubboResponse]{}).
		SetSpanNameExtractor(&rpc.RpcSpanNameExtractor[dubboRequest]{Getter: clientGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[dubboRequest]{}).
		AddAttributesExtractor(&rpc.ClientRpcAttrsExtractor[dubboRequest, dubboResponse, dubboAttrsGetter]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.DUBBO_CLIENT_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddOperationListeners(rpc.RpcClientMetrics("dubbo.client")).
		BuildInstrumenter()
}

func BuildDubboServerInstrumenter() instrumenter.Instrumenter[dubboRequest, dubboResponse] {
	builder := instrumenter.Builder[dubboRequest, dubboResponse]{}
	serverGetter := dubboAttrsGetter{}
	return builder.Init().
		SetSpanStatusExtractor(&dubboStatusExtractor[dubboRequest, dubboResponse]{}).
		SetSpanNameExtractor(&rpc.RpcSpanNameExtractor[dubboRequest]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[dubboRequest]{}).
		AddAttributesExtractor(&rpc.ServerRpcAttrsExtractor[dubboRequest, dubboResponse, dubboAttrsGetter]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.DUBBO_SERVER_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddOperationListeners(rpc.RpcServerMetrics("dubbo.server")).
		BuildInstrumenter()
}
