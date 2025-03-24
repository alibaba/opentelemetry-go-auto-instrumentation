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
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	_ "unsafe"
)

func publishOnEnter(call api.CallContext,
	ch *amqp.Channel,
	msg interface{},
) {
	if reflect.TypeOf(msg).String() != "*amqp091.basicPublish" {
		return
	}
	if !reflect.ValueOf(msg).Elem().IsValid() {
		return
	}
	request := RabbitRequest{
		operationName: "publish",
	}
	values := reflect.ValueOf(msg).Elem()

	if values.FieldByName("Exchange").IsValid() {
		request.exchange = values.FieldByName("Exchange").String()
	}
	if values.FieldByName("RoutingKey").IsValid() {
		request.routingKey = values.FieldByName("RoutingKey").String()
	}
	if values.FieldByName("Body").IsValid() {
		request.bodySize = int64(len(values.FieldByName("Body").Bytes()))
	}
	request.destinationName = request.exchange + ":" + request.routingKey
	ctx := context.Background()
	var attributes []attribute.KeyValue
	attributes = append(attributes,
		semconv.MessagingRabbitmqDestinationRoutingKey(request.routingKey), attribute.KeyValue{
			Key:   semconv.MessagingOperationTypeKey,
			Value: attribute.StringValue(request.operationName),
		},
	)
	if !values.FieldByName("Properties").IsValid() {
		RCtx := RabbitMQPublishEnabler.Start(ctx, &request, trace.WithAttributes(attributes...))
		data := make(map[string]interface{})
		data["ctx"] = RCtx
		data["rabbitMQ_publish_request"] = request
		call.SetData(data)
		return
	}
	if values.FieldByName("Properties").FieldByName("DeliveryMode").IsValid() {
		attributes = append(attributes,
			attribute.KeyValue{
				Key:   "messaging.rabbitmq.message.delivery_mode",
				Value: attribute.IntValue(int(values.FieldByName("Properties").FieldByName("DeliveryMode").Uint())),
			},
		)
	}
	if values.FieldByName("Properties").FieldByName("MessageId").IsValid() {
		request.messageId = values.FieldByName("Properties").FieldByName("MessageId").String()
	}
	var conversationIDValid bool
	var conversationID string
	if values.FieldByName("Properties").FieldByName("CorrelationId").IsValid() {
		conversationIDValid = true
		conversationID = values.FieldByName("Properties").FieldByName("CorrelationId").String()
		request.conversationID = conversationID
	}
	RCtx := RabbitMQPublishEnabler.Start(ctx, &request, trace.WithAttributes(attributes...))
	data := make(map[string]interface{})
	data["ctx"] = RCtx
	data["rabbitMQ_publish_request"] = request
	call.SetData(data)

	if conversationIDValid && request.conversationID != "" && conversationID == "" {
		reflect.ValueOf(msg).Elem().FieldByName("Properties").FieldByName("CorrelationId").SetString(request.conversationID)
	}

}
func publishOnExit(call api.CallContext, err error) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	request, ok := data["rabbitMQ_publish_request"].(RabbitRequest)
	if !ok {
		return
	}
	RabbitMQPublishEnabler.End(ctx, &request, nil, err)
}
