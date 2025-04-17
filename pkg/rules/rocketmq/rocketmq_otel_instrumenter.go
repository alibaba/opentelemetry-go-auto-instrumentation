package rocketmq

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/message"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
	"strconv"
)

// 埋点开关控制器
var rocketmqEnabler = instrumenter.NewDefaultInstrumentEnabler()

// 缓存Instrumenter实例，避免重复创建
var (
	producerInstrumenter      = buildRocketMQProducerInstrumenter()
	singleProcessInstrumenter = buildRocketMQConsumerInstrumenter(message.PROCESS, false)
	batchProcessInstrumenter  = buildRocketMQConsumerInstrumenter(message.PROCESS, true)
	receiveInstrumenter       = buildRocketMQConsumerReceiveInstrumenter()
)

// 将SendStatus转换为字符串
func sendStatusStr(status primitive.SendStatus) string {
	switch status {
	case primitive.SendOK:
		return "SEND_OK"
	case primitive.SendFlushDiskTimeout:
		return "SEND_FLUSH_DISK_TIMEOUT"
	case primitive.SendFlushSlaveTimeout:
		return "SEND_SENDFLUSH_SLAVE_TIMEOUT"
	case primitive.SendSlaveNotAvailable:
		return "SEND_SLAVE_NOT_AVAILABLE"
	default:
		return "UNKNOWN"
	}
}

// OpenTelemetry传播器载体
type rocketmqProducerCarrier struct {
	msg *primitive.Message
}

func (c rocketmqProducerCarrier) Get(key string) string {
	return c.msg.GetProperty(key)
}

func (c rocketmqProducerCarrier) Set(key, value string) {
	c.msg.WithProperty(key, value)
}

func (c rocketmqProducerCarrier) Keys() []string {
	if c.msg.GetProperties() == nil {
		return nil
	}
	keys := make([]string, 0, len(c.msg.GetProperties()))
	for k := range c.msg.GetProperties() {
		keys = append(keys, k)
	}

	return keys
}

// OpenTelemetry传播器载体
type rocketmqConsumerCarrier struct {
	msg *primitive.MessageExt
}

func (c rocketmqConsumerCarrier) Get(key string) string {
	return c.msg.GetProperty(key)
}

func (c rocketmqConsumerCarrier) Set(key, value string) {
	c.msg.WithProperty(key, value)
}

func (c rocketmqConsumerCarrier) Keys() []string {
	if c.msg.GetProperties() == nil {
		return nil
	}
	keys := make([]string, 0, len(c.msg.GetProperties()))
	for k := range c.msg.GetProperties() {
		keys = append(keys, k)
	}

	return keys
}

// 状态提取器
type rocketmqProducerStatusExtractor struct {
}

func (e *rocketmqProducerStatusExtractor) Extract(span trace.Span, req rocketmqProducerReq, res rocketmqProducerRes, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

type rocketmqConsumerStatusExtractor struct{}

func (e *rocketmqConsumerStatusExtractor) Extract(span trace.Span, req rocketmqConsumerReq, res rocketmqConsumerRes, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else if res.consumeResult != consumer.ConsumeSuccess {
		span.SetStatus(codes.Error, fmt.Sprintf("Consume failed with result: %v", res.consumeResult))
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

// rocketMQMessageProducerAttrsGetter 实现RocketMQ消息的属性获取接口
type rocketMQMessageProducerAttrsGetter struct{}

func (r rocketMQMessageProducerAttrsGetter) IsAnonymousDestination(request rocketmqProducerReq) bool {
	return false
}

func (r rocketMQMessageProducerAttrsGetter) GetDestinationPartitionId(request rocketmqProducerReq) string {
	if request.message == nil {
		return ""
	}
	return request.message.GetShardingKey()
}

// GetSystem 获取消息系统名称
func (r rocketMQMessageProducerAttrsGetter) GetSystem(request rocketmqProducerReq) string {
	return "rocketmq"
}

// GetDestination 获取目标主题
func (r rocketMQMessageProducerAttrsGetter) GetDestination(request rocketmqProducerReq) string {
	if request.message == nil {
		return ""
	}
	return request.message.Topic
}

// GetDestinationTemplate 获取目标主题模板
func (r rocketMQMessageProducerAttrsGetter) GetDestinationTemplate(request rocketmqProducerReq) string {
	return ""
}

// IsTemporaryDestination 判断是否是临时目标
func (r rocketMQMessageProducerAttrsGetter) IsTemporaryDestination(request rocketmqProducerReq) bool {
	return false
}

// isAnonymousDestination 判断是否是匿名目标
func (r rocketMQMessageProducerAttrsGetter) isAnonymousDestination(request rocketmqProducerReq) bool {
	return false
}

// GetConversationId 获取对话ID
func (r rocketMQMessageProducerAttrsGetter) GetConversationId(request rocketmqProducerReq) string {
	return ""
}

// GetMessageBodySize 获取消息体大小
func (r rocketMQMessageProducerAttrsGetter) GetMessageBodySize(request rocketmqProducerReq) int64 {
	if request.message == nil {
		return 0
	}
	if request.message.Batch && request.message.Compress {
		return int64(len(request.message.CompressedBody))
	}
	return int64(len(request.message.Body))
}

// GetMessageEnvelopSize 获取消息信封大小
func (r rocketMQMessageProducerAttrsGetter) GetMessageEnvelopSize(request rocketmqProducerReq) int64 {
	return 0
}

// GetMessageId 获取消息ID
func (r rocketMQMessageProducerAttrsGetter) GetMessageId(request rocketmqProducerReq, response rocketmqProducerRes) string {
	if response.result != nil {
		return response.result.MsgID
	}
	return ""
}

// GetClientId 获取客户端ID
func (r rocketMQMessageProducerAttrsGetter) GetClientId(request rocketmqProducerReq) string {
	return request.clientID
}

// GetBatchMessageCount 获取批处理消息数量
func (r rocketMQMessageProducerAttrsGetter) GetBatchMessageCount(request rocketmqProducerReq, response rocketmqProducerRes) int64 {
	return 1 // RocketMQ默认每次发送单条消息
}

// GetMessageHeader 获取消息头
func (r rocketMQMessageProducerAttrsGetter) GetMessageHeader(request rocketmqProducerReq, name string) []string {
	if request.message == nil {
		return nil
	}
	value := request.message.GetProperty(name)
	if value == "" {
		return nil
	}
	return []string{value}
}

// rocketMQMessageConsumerAttrsGetter 实现RocketMQ消费者的属性获取接口
type rocketMQMessageConsumerAttrsGetter struct{}

func (r rocketMQMessageConsumerAttrsGetter) IsAnonymousDestination(request rocketmqConsumerReq) bool {
	return false
}

func (r rocketMQMessageConsumerAttrsGetter) GetDestinationPartitionId(request rocketmqConsumerReq) string {
	return ""
}

// GetSystem 获取消息系统名称
func (r rocketMQMessageConsumerAttrsGetter) GetSystem(request rocketmqConsumerReq) string {
	return "rocketmq"
}

// GetDestination 获取目标主题
func (r rocketMQMessageConsumerAttrsGetter) GetDestination(request rocketmqConsumerReq) string {
	if request.messages == nil {
		return ""
	}
	return request.messages.Topic
}

// GetDestinationTemplate 获取目标主题模板
func (r rocketMQMessageConsumerAttrsGetter) GetDestinationTemplate(request rocketmqConsumerReq) string {
	return ""
}

// IsTemporaryDestination 判断是否是临时目标
func (r rocketMQMessageConsumerAttrsGetter) IsTemporaryDestination(request rocketmqConsumerReq) bool {
	return false
}

// isAnonymousDestination 判断是否是匿名目标
func (r rocketMQMessageConsumerAttrsGetter) isAnonymousDestination(request rocketmqConsumerReq) bool {
	return false
}

// GetConversationId 获取对话ID
func (r rocketMQMessageConsumerAttrsGetter) GetConversationId(request rocketmqConsumerReq) string {
	return ""
}

// GetMessageBodySize 获取消息体大小
func (r rocketMQMessageConsumerAttrsGetter) GetMessageBodySize(request rocketmqConsumerReq) int64 {
	if request.messages == nil {
		return 0
	}
	return int64(len(request.messages.Body))
}

// GetMessageEnvelopSize 获取消息信封大小
func (r rocketMQMessageConsumerAttrsGetter) GetMessageEnvelopSize(request rocketmqConsumerReq) int64 {
	return 0
}

// GetMessageId 获取消息ID
func (r rocketMQMessageConsumerAttrsGetter) GetMessageId(request rocketmqConsumerReq, response rocketmqConsumerRes) string {
	if request.messages == nil {
		return ""
	}
	return request.messages.MsgId
}

// GetClientId 获取客户端ID
func (r rocketMQMessageConsumerAttrsGetter) GetClientId(request rocketmqConsumerReq) string {
	return ""
}

// GetBatchMessageCount 获取批处理消息数量
func (r rocketMQMessageConsumerAttrsGetter) GetBatchMessageCount(request rocketmqConsumerReq, response rocketmqConsumerRes) int64 {
	return 1
}

// GetMessageHeader 获取消息头
func (r rocketMQMessageConsumerAttrsGetter) GetMessageHeader(request rocketmqConsumerReq, name string) []string {
	if request.messages == nil {
		return nil
	}
	value := request.messages.GetProperty(name)
	if value == "" {
		return nil
	}
	return []string{value}
}

type rocketmqConsumerProcessAttributesExtractor struct{}

func (e *rocketmqConsumerProcessAttributesExtractor) OnStart(attributes []attribute.KeyValue, parentContext context.Context, req rocketmqConsumerReq) ([]attribute.KeyValue, context.Context) {
	if req.messages == nil {
		return attributes, parentContext
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingRocketmqMessageTagKey.String(req.messages.GetTags()),
		semconv.MessagingRocketmqMessageKeysKey.String(req.messages.GetKeys()),
	}
	return append(attributes, attrs...), parentContext
}

func (e *rocketmqConsumerProcessAttributesExtractor) OnEnd(attributes []attribute.KeyValue, ctx context.Context, req rocketmqConsumerReq, res rocketmqConsumerRes, err error) ([]attribute.KeyValue, context.Context) {
	if err != nil {
		attributes = append(attributes, attribute.String("error", err.Error()))
	}

	if req.addr != "" {
		attributes = append(attributes, attribute.String("messaging.rocketmq.broker_address", req.addr))
	}

	if req.messages == nil {
		return attributes, ctx
	}

	attributes = append(attributes,
		attribute.String("messaging.rocketmq.queue_offset", strconv.FormatInt(req.messages.QueueOffset, 10)),
	)
	if req.messages.TransactionId != "" {
		attributes = append(attributes, attribute.String("messaging.rocketmq.transaction_id", req.messages.TransactionId))
	}
	if req.messages.Queue != nil {
		attributes = append(attributes,
			attribute.Int("messaging.rocketmq.queue_id", req.messages.Queue.QueueId),
			attribute.String("messaging.rocketmq.broker_name", req.messages.Queue.BrokerName),
		)
	}

	return attributes, ctx
}

// 属性提取器
type rocketmqProducerAttributesExtractor struct {
}

func (m *rocketmqProducerAttributesExtractor) OnStart(attributes []attribute.KeyValue, parentContext context.Context, req rocketmqProducerReq) ([]attribute.KeyValue, context.Context) {
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemRocketmq,
		semconv.MessagingDestinationNameKey.String(req.message.Topic),
		semconv.MessagingRocketmqMessageTagKey.String(req.message.GetTags()),
		semconv.MessagingRocketmqMessageKeysKey.String(req.message.GetKeys()),
	}

	return append(attributes, attrs...), parentContext
}

type rocketmqConsumerSpanNameExtractor struct {
}

func (e *rocketmqConsumerSpanNameExtractor) Extract(req any) string {
	return "multiple_sources receive"
}

func (e *rocketmqProducerAttributesExtractor) OnEnd(attributes []attribute.KeyValue, ctx context.Context, req rocketmqProducerReq, res rocketmqProducerRes, err error) ([]attribute.KeyValue, context.Context) {
	if err != nil {
		attributes = append(attributes, attribute.String("error", err.Error()))
	}
	if res.addr != "" {
		attributes = append(attributes, attribute.String("messaging.rocketmq.broker_address", res.addr))
	}

	if res.result != nil {
		attributes = append(attributes,
			attribute.String("messaging.rocketmq.queue_offset", strconv.FormatInt(res.result.QueueOffset, 10)),
			attribute.String("messaging.rocketmq.send_result", sendStatusStr(res.result.Status)),
			semconv.MessagingClientID(req.clientID),
		)
		if res.result.TransactionID != "" {
			attributes = append(attributes, attribute.String("messaging.rocketmq.transaction_id", res.result.TransactionID))
		}
		if res.result.MessageQueue != nil {
			attributes = append(attributes,
				attribute.Int("messaging.rocketmq.queue_id", res.result.MessageQueue.QueueId),
				attribute.String("messaging.rocketmq.broker_name", res.result.MessageQueue.BrokerName),
			)
		}
	}
	attributes = append(attributes,
		semconv.MessagingClientID(req.clientID),
	)
	return attributes, ctx
}

// 构建Producer埋点器
func buildRocketMQProducerInstrumenter() instrumenter.Instrumenter[rocketmqProducerReq, rocketmqProducerRes] {
	builder := instrumenter.Builder[rocketmqProducerReq, rocketmqProducerRes]{}
	return builder.Init().
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.ROCKETMQGO_PRODUCER_SCOPE_NAME,
			Version: version.Tag,
		}).
		SetSpanNameExtractor(&message.MessageSpanNameExtractor[rocketmqProducerReq, rocketmqProducerRes]{Getter: rocketMQMessageProducerAttrsGetter{}, OperationName: message.PUBLISH}).
		SetSpanKindExtractor(&instrumenter.AlwaysProducerExtractor[rocketmqProducerReq]{}).
		AddAttributesExtractor(&rocketmqProducerAttributesExtractor{}).
		AddAttributesExtractor(&message.MessageAttrsExtractor[rocketmqProducerReq, rocketmqProducerRes, rocketMQMessageProducerAttrsGetter]{Operation: message.PUBLISH}).
		SetSpanStatusExtractor(&rocketmqProducerStatusExtractor{}).
		BuildPropagatingToDownstreamInstrumenter(
			func(req rocketmqProducerReq) propagation.TextMapCarrier {
				return rocketmqProducerCarrier{msg: req.message}
			},
			otel.GetTextMapPropagator(),
		)
}

// 构建消费埋点器
func buildRocketMQConsumerReceiveInstrumenter() instrumenter.Instrumenter[any, any] {
	builder := instrumenter.Builder[any, any]{}
	return builder.Init().
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.ROCKETMQGO_CONSUMER_SCOPE_NAME,
			Version: version.Tag,
		}).
		SetSpanNameExtractor(&rocketmqConsumerSpanNameExtractor{}).
		SetSpanKindExtractor(&instrumenter.AlwaysConsumerExtractor[any]{}).
		BuildInstrumenter()
}

// 统一构建RocketMQ消费者Instrumenter的方法
func buildRocketMQConsumerInstrumenter(operationName message.MessageOperation, isBatch bool) instrumenter.Instrumenter[rocketmqConsumerReq, rocketmqConsumerRes] {
	builder := instrumenter.Builder[rocketmqConsumerReq, rocketmqConsumerRes]{}
	buildInstrumenter := builder.Init().
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.ROCKETMQGO_CONSUMER_SCOPE_NAME,
			Version: version.Tag,
		}).
		SetSpanNameExtractor(&message.MessageSpanNameExtractor[rocketmqConsumerReq, rocketmqConsumerRes]{
			Getter:        rocketMQMessageConsumerAttrsGetter{},
			OperationName: operationName,
		}).
		SetSpanKindExtractor(&instrumenter.AlwaysConsumerExtractor[rocketmqConsumerReq]{}).
		AddAttributesExtractor(&rocketmqConsumerProcessAttributesExtractor{}).
		AddAttributesExtractor(&message.MessageAttrsExtractor[rocketmqConsumerReq, rocketmqConsumerRes, rocketMQMessageConsumerAttrsGetter]{
			Operation: operationName,
		})

	if !isBatch {
		return buildInstrumenter.SetSpanStatusExtractor(&rocketmqConsumerStatusExtractor{}).
			BuildPropagatingFromUpstreamInstrumenter(
				func(req rocketmqConsumerReq) propagation.TextMapCarrier {
					return rocketmqConsumerCarrier{msg: req.messages}
				},
				otel.GetTextMapPropagator(),
			)
	} else {
		return buildInstrumenter.BuildInstrumenter()
	}
}

func Start(parentContext context.Context, msgs []*primitive.MessageExt, addr string) context.Context {
	if len(msgs) == 1 {
		var messages *primitive.MessageExt
		if len(msgs) != 0 || msgs[0] != nil {
			messages = msgs[0]
		}
		return singleProcessInstrumenter.Start(parentContext, rocketmqConsumerReq{addr: addr, messages: messages})
	} else {
		var attributes []attribute.KeyValue
		attributes = append(attributes,
			semconv.MessagingSystemRocketmq,
			attribute.KeyValue{
				Key:   semconv.MessagingOperationTypeKey,
				Value: attribute.StringValue(string(message.RECEIVE)),
			},
		)
		rootContext := receiveInstrumenter.Start(parentContext, nil, trace.WithAttributes(attributes...))
		for _, msg := range msgs {
			createChildSpan(rootContext, msg, addr)
		}
		return rootContext
	}
}

func createChildSpan(ctx context.Context, msg *primitive.MessageExt, addr string) {
	context := batchProcessInstrumenter.Start(ctx, rocketmqConsumerReq{addr: addr, messages: msg}, trace.WithLinks(trace.Link{SpanContext: trace.SpanContextFromContext(otel.GetTextMapPropagator().Extract(ctx, rocketmqConsumerCarrier{msg: msg}))}))
	batchProcessInstrumenter.End(context, rocketmqConsumerReq{addr: addr, messages: msg}, rocketmqConsumerRes{}, nil)
}

func End(ctx context.Context, request []*primitive.MessageExt, response consumer.ConsumeResult, err error) {
	if len(request) == 1 {
		var messages *primitive.MessageExt
		if len(request) != 0 || request[0] != nil {
			messages = request[0]
		}
		singleProcessInstrumenter.End(ctx, rocketmqConsumerReq{messages: messages}, rocketmqConsumerRes{response}, err)
	} else {
		receiveInstrumenter.End(ctx, request, response, err)
	}
}
