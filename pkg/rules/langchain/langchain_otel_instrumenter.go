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

package langchain

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type LcExperimentalSpanNameExtractor struct {
}

func (l LcExperimentalSpanNameExtractor) Extract(request langChainRequest) string {
	if request.operationName == "" {
		return "unknown module name"
	}
	return request.operationName
}

type LcExperimentalAttributeExtractor struct {
}

func (l LcExperimentalAttributeExtractor) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request langChainRequest) ([]attribute.KeyValue, context.Context) {
	attributes = append(attributes, attribute.KeyValue{
		Key:   "module-name",
		Value: attribute.StringValue(request.operationName),
	}, attribute.KeyValue{
		Key:   "system",
		Value: attribute.StringValue("langChain"),
	})
	if request.input != nil {
		var val attribute.Value
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
					Key:   attribute.Key("input." + k),
					Value: val,
				})
			}
			val = attribute.Value{}
		}

	}
	return attributes, parentContext
}

func (l LcExperimentalAttributeExtractor) OnEnd(attributes []attribute.KeyValue, context context.Context, request langChainRequest, response any, err error) ([]attribute.KeyValue, context.Context) {
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
					Key:   attribute.Key("output." + k),
					Value: val,
				})
			}
			val = attribute.Value{}
		}

	}
	return attributes, context
}

func BuildLangChainInternalInstrumenter() instrumenter.Instrumenter[langChainRequest, any] {
	builder := instrumenter.Builder[langChainRequest, any]{}
	return builder.Init().SetSpanNameExtractor(&LcExperimentalSpanNameExtractor{}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[langChainRequest]{}).
		AddAttributesExtractor(&LcExperimentalAttributeExtractor{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.LANGCHAIN_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
