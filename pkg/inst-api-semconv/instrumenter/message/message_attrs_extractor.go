package message

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
)

type MessageOperation string

const PUBLISH MessageOperation = "publish"
const RECEIVE MessageOperation = "receive"
const PROCESS MessageOperation = "process"

const messaging_batch_message_count = attribute.Key("messaging.batch.message_count")
const messaging_client_id = attribute.Key("messaging.client_id")
const messaging_destination_anoymous = attribute.Key("messaging.destination.anonymous")
const messaging_destination_name = attribute.Key("messaging.destination.name")
const messaging_destination_template = attribute.Key("messaging.destination.template")
const messaging_destination_temporary = attribute.Key("messaging.destination.temporary")
const messaging_message_body_size = attribute.Key("messaging.message.body.size")
const messaging_message_conversation_id = attribute.Key("messaging.message.conversation_id")
const messaging_message_envelope_size = attribute.Key("messaging.message.envelope.size")
const messaging_message_id = attribute.Key("messaging.message.id")
const messaging_operation = attribute.Key("messaging.operation")
const messaging_system = attribute.Key("messaging.system")

type MessageAttrsExtractor[REQUEST any, RESPONSE any, GETTER MessageAttrsGetter[REQUEST, RESPONSE]] struct {
	getter    GETTER
	operation MessageOperation
}

func (m *MessageAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	switch m.operation {
	case PUBLISH:
		return utils.PRODUCER_KEY
	case RECEIVE:
		return utils.CONSUMER_RECEIVE_KEY
	case PROCESS:
		return utils.CONSUMER_PROCESS_KEY
	}
	panic("Operation" + m.operation + "not supported")
}

func (m *MessageAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	messageAttrSystem := m.getter.GetSystem(request)
	isTemporaryDestination := m.getter.IsTemporaryDestination(request)
	if isTemporaryDestination {
		attributes = append(attributes, attribute.KeyValue{
			Key:   messaging_destination_temporary,
			Value: attribute.BoolValue(true),
		}, attribute.KeyValue{
			Key:   messaging_destination_name,
			Value: attribute.StringValue("(temporary)"),
		})
	} else {
		attributes = append(attributes, attribute.KeyValue{
			Key:   messaging_destination_name,
			Value: attribute.StringValue(m.getter.GetDestination(request)),
		}, attribute.KeyValue{
			Key:   messaging_destination_template,
			Value: attribute.StringValue(m.getter.GetDestinationTemplate(request)),
		})
	}
	isAnonymousDestination := m.getter.isAnonymousDestination(request)
	if isAnonymousDestination {
		attributes = append(attributes, attribute.KeyValue{
			Key:   messaging_destination_anoymous,
			Value: attribute.BoolValue(true),
		})
	}
	attributes = append(attributes, attribute.KeyValue{
		Key:   messaging_message_conversation_id,
		Value: attribute.StringValue(m.getter.GetConversationId(request)),
	}, attribute.KeyValue{
		Key:   messaging_message_body_size,
		Value: attribute.Int64Value(m.getter.GetMessageBodySize(request)),
	}, attribute.KeyValue{
		Key:   messaging_message_envelope_size,
		Value: attribute.Int64Value(m.getter.GetMessageEnvelopSize(request)),
	}, attribute.KeyValue{
		Key:   messaging_client_id,
		Value: attribute.StringValue(m.getter.GetClientId(request)),
	}, attribute.KeyValue{
		Key:   messaging_operation,
		Value: attribute.StringValue(string(m.operation)),
	}, attribute.KeyValue{
		Key:   messaging_system,
		Value: attribute.StringValue(messageAttrSystem),
	})
	return attributes
}

func (m *MessageAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   messaging_message_id,
		Value: attribute.StringValue(m.getter.GetMessageId(request, response)),
	}, attribute.KeyValue{
		Key:   messaging_batch_message_count,
		Value: attribute.Int64Value(m.getter.GetBatchMessageCount(request, response)),
	})
	// TODO: add custom captured headers attributes
	return attributes
}
