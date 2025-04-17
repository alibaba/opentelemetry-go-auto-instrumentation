package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// 定义请求响应类型
type rocketmqProducerReq struct {
	//ctx     context.Context
	clientID string
	message  *primitive.Message
}

type rocketmqProducerRes struct {
	addr   string
	result *primitive.SendResult
}

// 消费处理请求类型
type rocketmqConsumerReq struct {
	addr     string
	messages *primitive.MessageExt
}

// 消费处理响应类型
type rocketmqConsumerRes struct {
	consumeResult consumer.ConsumeResult
}
