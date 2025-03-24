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

package amqp091go

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/message"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type RabbitMQGetter struct {
}

var _ message.MessageAttrsGetter[*RabbitRequest, any] = RabbitMQGetter{}

func (RabbitMQGetter) GetSystem(request *RabbitRequest) string {
	return "rabbitmq"
}

func (RabbitMQGetter) GetDestination(request *RabbitRequest) string {
	return request.destinationName
}

func (RabbitMQGetter) GetDestinationTemplate(request *RabbitRequest) string {
	return ""
}

func (RabbitMQGetter) IsTemporaryDestination(request *RabbitRequest) bool {
	return false
}

func (RabbitMQGetter) IsAnonymousDestination(request *RabbitRequest) bool {
	return false
}

func (RabbitMQGetter) GetConversationId(request *RabbitRequest) string {
	return request.conversationID
}

func (RabbitMQGetter) GetMessageBodySize(request *RabbitRequest) int64 {
	return request.bodySize
}

func (RabbitMQGetter) GetMessageEnvelopSize(request *RabbitRequest) int64 {
	return 0
}

func (RabbitMQGetter) GetMessageId(request *RabbitRequest, response any) string {
	return request.messageId
}

func (RabbitMQGetter) GetClientId(request *RabbitRequest) string {
	return ""
}

func (RabbitMQGetter) GetBatchMessageCount(request *RabbitRequest, response any) int64 {
	return 0
}

func (RabbitMQGetter) GetMessageHeader(request *RabbitRequest, name string) []string {
	return []string{}
}

func BuildRabbitMQConsumeOtelInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[*RabbitRequest, any] {
	builder := instrumenter.Builder[*RabbitRequest, any]{}
	return builder.Init().SetSpanNameExtractor(&message.MessageSpanNameExtractor[*RabbitRequest, any]{Getter: RabbitMQGetter{}, OperationName: message.RECEIVE}).
		SetSpanKindExtractor(&instrumenter.AlwaysConsumerExtractor[*RabbitRequest]{}).
		AddAttributesExtractor(&message.MessageAttrsExtractor[*RabbitRequest, any, RabbitMQGetter]{Operation: message.RECEIVE}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.AMQP091GO_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildPropagatingFromUpstreamInstrumenter(func(n *RabbitRequest) propagation.TextMapCarrier {
			return n
		}, otel.GetTextMapPropagator())
}
func BuildRabbitMQPublishOtelInstrumenter() *instrumenter.PropagatingToDownstreamInstrumenter[*RabbitRequest, any] {
	builder := instrumenter.Builder[*RabbitRequest, any]{}
	return builder.Init().SetSpanNameExtractor(&message.MessageSpanNameExtractor[*RabbitRequest, any]{Getter: RabbitMQGetter{}, OperationName: message.PUBLISH}).
		SetSpanKindExtractor(&instrumenter.AlwaysProducerExtractor[*RabbitRequest]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.AMQP091GO_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddAttributesExtractor(&message.MessageAttrsExtractor[*RabbitRequest, any, RabbitMQGetter]{Operation: message.PUBLISH}).
		BuildPropagatingToDownstreamInstrumenter(func(n *RabbitRequest) propagation.TextMapCarrier {
			return n
		}, otel.GetTextMapPropagator())

}
