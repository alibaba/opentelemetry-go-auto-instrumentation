package eino

import (
	"context"
	"fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/ai"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type einoCommonAttrsGetter struct {
}

var _ ai.CommonAttrsGetter[einoRequest, einoResponse] = einoCommonAttrsGetter{}

func (einoCommonAttrsGetter) GetAIOperationName(request einoRequest) string {
	return request.operationName
}
func (einoCommonAttrsGetter) GetAISystem(request einoRequest) string {
	return "eino"
}

type LExperimentalAttributeExtractor struct {
	Base ai.AICommonAttrsExtractor[einoRequest, any, einoCommonAttrsGetter]
}

func (l LExperimentalAttributeExtractor) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request einoRequest) ([]attribute.KeyValue, context.Context) {
	attributes, parentContext = l.Base.OnStart(attributes, parentContext, request)
	var val attribute.Value
	if request.input != nil {
		for k, v := range request.input {
			val = attribute.StringValue(fmt.Sprintf("%#v", v))
			attributes = append(attributes, attribute.KeyValue{
				Key:   attribute.Key(fmt.Sprintf("gen_ai.%s.%s", request.operationName, k)),
				Value: val,
			})
			val = attribute.Value{}
		}
	}
	return attributes, parentContext
}

func (l LExperimentalAttributeExtractor) OnEnd(attributes []attribute.KeyValue, context context.Context, request einoRequest, response einoResponse, err error) ([]attribute.KeyValue, context.Context) {
	attributes, context = l.Base.OnEnd(attributes, context, request, response, err)
	if response.output != nil {
		var val attribute.Value
		for k, v := range response.output {
			val = attribute.StringValue(fmt.Sprintf("%#v", v))
			attributes = append(attributes, attribute.KeyValue{
				Key:   attribute.Key(fmt.Sprintf("gen_ai.%s.%s", request.operationName, k)),
				Value: val,
			})
			val = attribute.Value{}
		}
	}
	return attributes, context
}

func BuildEinoCommonInstrumenter() instrumenter.Instrumenter[einoRequest, einoResponse] {
	builder := instrumenter.Builder[einoRequest, einoResponse]{}
	return builder.Init().SetSpanNameExtractor(&ai.AISpanNameExtractor[einoRequest, einoResponse]{Getter: einoCommonAttrsGetter{}}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[einoRequest]{}).
		AddAttributesExtractor(&LExperimentalAttributeExtractor{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.EINO_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
