package mcp

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/ai"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"github.com/mark3labs/mcp-go/mcp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type aiCommonRequest struct {
}

func (aiCommonRequest) GetAIOperationName(request mcpServerRequest) string {
	return request.operationName
}
func (aiCommonRequest) GetAISystem(request mcpServerRequest) string {
	return request.system
}

type LExperimentalAttributeExtractor struct {
	Base ai.AICommonAttrsExtractor[mcpServerRequest, any, aiCommonRequest]
}

func (l LExperimentalAttributeExtractor) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request mcpServerRequest) ([]attribute.KeyValue, context.Context) {
	attributes, parentContext = l.Base.OnStart(attributes, parentContext, request)
	var val attribute.Value
	if request.methodType == string(mcp.MethodToolsCall) {
		attributes = append(attributes, attribute.KeyValue{
			Key:   "gen_ai.tool.name",
			Value: attribute.StringValue(request.methodName),
		}, attribute.KeyValue{
			Key:   "gen_ai.tool.call.id",
			Value: attribute.StringValue(request.CallId),
		})
	}
	if request.input != nil {
		for k, v := range request.input {
			switch v.(type) {
			case string:
				val = attribute.StringValue(v.(string))
			case int:
				val = attribute.IntValue(v.(int))
			case int64:
				val = attribute.Int64Value(v.(int64))
			case float64:
				val = attribute.Float64Value(v.(float64))
			case bool:
				val = attribute.BoolValue(v.(bool))
			default:
				val = attribute.StringValue(fmt.Sprintf("%#v", v))
			}
			if val.Type() > 0 {
				attributes = append(attributes, attribute.KeyValue{
					Key:   attribute.Key("gen_ai.other_input." + k),
					Value: val,
				})
			}
			val = attribute.Value{}
		}
	}

	return attributes, parentContext
}

func (l LExperimentalAttributeExtractor) OnEnd(attributes []attribute.KeyValue, context context.Context, request mcpServerRequest, response any, err error) ([]attribute.KeyValue, context.Context) {
	attributes, context = l.Base.OnEnd(attributes, context, request, response, err)
	if request.output != nil {
		var val attribute.Value
		for k, v := range request.output {
			switch v.(type) {
			case string:
				val = attribute.StringValue(v.(string))
			case int:
				val = attribute.IntValue(v.(int))
			case int64:
				val = attribute.Int64Value(v.(int64))
			case float64:
				val = attribute.Float64Value(v.(float64))
			case bool:
				val = attribute.BoolValue(v.(bool))
			default:
				val = attribute.StringValue(fmt.Sprintf("%#v", v))
			}
			if val.Type() > 0 {
				attributes = append(attributes, attribute.KeyValue{
					Key:   attribute.Key("gen_ai.other_output." + k),
					Value: val,
				})
			}
			val = attribute.Value{}
		}

	}
	return attributes, context
}

func BuildCommonOtelInstrumenter() instrumenter.Instrumenter[mcpServerRequest, any] {
	builder := instrumenter.Builder[mcpServerRequest, any]{}
	return builder.Init().SetSpanNameExtractor(&ai.AISpanNameExtractor[mcpServerRequest, any]{Getter: aiCommonRequest{}}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[mcpServerRequest]{}).
		AddAttributesExtractor(&LExperimentalAttributeExtractor{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.MCP_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
