// Copyright (c) 2025 Alibaba Group Holding Ltd.
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
