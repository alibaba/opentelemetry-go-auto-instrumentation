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

package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"
)

const (
	topicName = "test_topic"
	GroupName = "test_group"
)

func initRocketMQ() rocketmq.Producer {
	// 从环境变量获取NameServer地址
	nameSrvAddr := os.Getenv("NAMESRV_ADDR")
	if nameSrvAddr == "" {
		nameSrvAddr = "127.0.0.1:9876"
	}

	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{nameSrvAddr}),
		producer.WithGroupName(GroupName),
		producer.WithRetry(2),
	)
	if err != nil {
		panic(err)
	}

	// 启动生产者
	err = p.Start()
	if err != nil {
		panic(err)
	}

	return p
}

// 创建测试主题
func initTopic() {
	nameSrvAddr := os.Getenv("NAMESRV_ADDR")
	if nameSrvAddr == "" {
		nameSrvAddr = "127.0.0.1:9876"
	}
	brokerAddr := os.Getenv("BROKER_ADDR")
	if brokerAddr == "" {
		brokerAddr = "127.0.0.1:10911"
	}
	// 创建admin客户端
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver([]string{nameSrvAddr})),
	)
	if err != nil {
		log.Printf("创建admin客户端失败: %s\n", err.Error())
		// 如果无法创建admin客户端，依赖RocketMQ自动创建主题功能
		return
	}
	defer adminClient.Close()

	// 创建主题
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = adminClient.CreateTopic(
		ctx,
		admin.WithTopicCreate(topicName),
		admin.WithBrokerAddrCreate(brokerAddr),
	)
	if err != nil {
		log.Printf("创建主题失败: %s\n", err.Error())
		log.Println("将依赖RocketMQ自动创建主题功能")
	} else {
		log.Printf("主题 %s 创建成功\n", topicName)
	}

	// 等待主题创建完成
	waitForTopicReadyByAdmin(adminClient, topicName, 1000, time.Second)

}

func waitForTopicReadyByAdmin(adminClient admin.Admin, topic string, maxRetries int, interval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		topics, err := adminClient.FetchAllTopicList(context.Background())
		if err != nil {
			log.Printf("获取 Topic 列表失败: %v", err)
		} else {
			for _, t := range topics.TopicList {
				if t == topic {
					log.Printf("Topic %s 已就绪", topic)
					return nil
				}
			}
		}
		log.Printf("第 %d 次检查: Topic %s 尚未就绪", i+1, topic)
		time.Sleep(interval)
	}
	return fmt.Errorf("等待 Topic %s 就绪超时", topic)
}

// 创建并配置消费者
func initConsumer() rocketmq.PushConsumer {
	// 从环境变量获取NameServer地址
	nameSrvAddr := os.Getenv("NAMESRV_ADDR")
	if nameSrvAddr == "" {
		nameSrvAddr = "127.0.0.1:9876"
	}

	// 创建消费者实例
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{nameSrvAddr}),
		consumer.WithGroupName(GroupName),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeMessageBatchMaxSize(10),
	)
	if err != nil {
		panic(err)
	}

	return c
}

// VerifyRocketMQAttributes 校验RocketMQProducer消息的Span属性
func VerifyRocketMQAttributes(producer tracetest.SpanStub, consumer tracetest.SpanStub, topic string, tag string) {
	verifier.VerifyMQPublishAttributes(producer, "", "", topic, "publish", topic, "rocketmq")
	verifier.VerifyMQConsumeAttributes(consumer, "", "", topic, "process", topic, "rocketmq")

	// 验证生产端和消费端trace span
	verifier.Assert(consumer.Parent.SpanID() == producer.SpanContext.SpanID(), "期望spanid为 %s, 实际为 %s", consumer.Parent.SpanID(), producer.SpanContext.SpanID())
	verifier.Assert(consumer.SpanContext.TraceID() == producer.SpanContext.TraceID(), "期望traceid为 %s, 实际为 %s", consumer.SpanContext.TraceID(), producer.SpanContext.TraceID())
}

// VerifyRocketMQProduceAttributes 校验RocketMQProducer消息的Span属性
func VerifyRocketMQProduceAttributes(span tracetest.SpanStub, topic string, tag string, key string, operationName string, expectedError bool) {
	// 验证基本的消息属性
	verifier.Assert(span.Name == topic+" "+operationName, "期望span名称为 %s, 实际为 %s", topic+" "+operationName, span.Name)

	// 验证标准消息属性
	actualSystem := verifier.GetAttribute(span.Attributes, "messaging.system").AsString()
	verifier.Assert(actualSystem == "rocketmq", "期望messaging.system为 %s, 实际为 %s", "rocketmq", actualSystem)

	actualDestination := verifier.GetAttribute(span.Attributes, "messaging.destination.name").AsString()
	verifier.Assert(actualDestination == topic, "期望messaging.destination.name为 %s, 实际为 %s", topic, actualDestination)

	actualOperation := verifier.GetAttribute(span.Attributes, "messaging.operation.name").AsString()
	verifier.Assert(actualOperation == operationName, "期望messaging.operation.name为 %s, 实际为 %s", operationName, actualOperation)

	// 验证消息体大小 (仅检查是否存在)
	bodySize := verifier.GetAttribute(span.Attributes, "messaging.message.body.size").AsInt64()
	verifier.Assert(bodySize > 0, "期望messaging.message.body.size大于0, 实际为 %d", bodySize)

	// 验证RocketMQ特有属性
	if tag != "" {
		actualTag := verifier.GetAttribute(span.Attributes, "messaging.rocketmq.message.tag").AsString()
		verifier.Assert(actualTag == tag, "期望messaging.rocketmq.message.tag为 %s, 实际为 %s", tag, actualTag)
	}

	if key != "" {
		actualKey := verifier.GetAttribute(span.Attributes, "messaging.rocketmq.message.keys").AsString()
		verifier.Assert(actualKey == key, "期望messaging.rocketmq.message.keys为 %s, 实际为 %s", key, actualKey)
	}

	// 验证Span种类
	verifier.Assert(span.SpanKind == trace.SpanKindProducer, "期望为生产者Span, 实际为 %d", span.SpanKind)

	// 验证错误状态
	if expectedError {
		verifier.Assert(span.Status.Code == codes.Error, "期望Span状态为Error, 实际为 %s", span.Status.Code)
		verifier.Assert(span.Status.Description != "", "期望Span错误描述不为空")
	} else {
		verifier.Assert(span.Status.Code != codes.Error, "期望Span状态不为Error, 实际为 %s", span.Status.Code)
	}
}

// VerifyRocketMQConsumeAttributes 校验RocketMQ消息的Span属性
func VerifyRocketMQConsumeAttributes(span tracetest.SpanStub, topic string, tag string, key string, operationName string, expectedError bool) {
	// 验证基本的消息属性
	verifier.Assert(span.Name == topic+" "+operationName, "期望span名称为 %s, 实际为 %s", topic+" "+operationName, span.Name)

	// 验证标准消息属性
	actualSystem := verifier.GetAttribute(span.Attributes, "messaging.system").AsString()
	verifier.Assert(actualSystem == "rocketmq", "期望messaging.system为 %s, 实际为 %s", "rocketmq", actualSystem)

	actualDestination := verifier.GetAttribute(span.Attributes, "messaging.destination.name").AsString()
	verifier.Assert(actualDestination == topic, "期望messaging.destination.name为 %s, 实际为 %s", topic, actualDestination)

	actualOperation := verifier.GetAttribute(span.Attributes, "messaging.operation.name").AsString()
	verifier.Assert(actualOperation == operationName, "期望messaging.operation.name为 %s, 实际为 %s", operationName, actualOperation)

	// 验证消息体大小 (仅检查是否存在)
	bodySize := verifier.GetAttribute(span.Attributes, "messaging.message.body.size").AsInt64()
	verifier.Assert(bodySize > 0, "期望messaging.message.body.size大于0, 实际为 %d", bodySize)

	// 验证RocketMQ特有属性
	if tag != "" {
		actualTag := verifier.GetAttribute(span.Attributes, "messaging.rocketmq.message.tag").AsString()
		verifier.Assert(actualTag == tag, "期望messaging.rocketmq.message.tag为 %s, 实际为 %s", tag, actualTag)
	}

	if key != "" {
		actualKey := verifier.GetAttribute(span.Attributes, "messaging.rocketmq.message.keys").AsString()
		verifier.Assert(actualKey == key, "期望messaging.rocketmq.message.keys为 %s, 实际为 %s", key, actualKey)
	}

	// 验证Span种类
	verifier.Assert(span.SpanKind == trace.SpanKindConsumer, "期望为消费者Span, 实际为 %d", span.SpanKind)

	// 验证错误状态
	if expectedError {
		verifier.Assert(span.Status.Code == codes.Error, "期望Span状态为Error, 实际为 %s", span.Status.Code)
		verifier.Assert(span.Status.Description != "", "期望Span错误描述不为空")
	} else {
		verifier.Assert(span.Status.Code != codes.Error, "期望Span状态不为Error, 实际为 %s", span.Status.Code)
	}
}

// VerifyRocketMQReceive 校验RocketMQ消息的Span属性
func VerifyRocketMQReceive(producter tracetest.SpanStub, receive tracetest.SpanStub, process tracetest.SpanStub) {
	// 验证基本的消息属性
	verifier.Assert(receive.Name == "multiple_sources receive", "期望span名称为 multiple_sources receive, 实际为 %s", receive.Name)
	// 验证标准消息属性
	actualSystem := verifier.GetAttribute(receive.Attributes, "messaging.system").AsString()
	verifier.Assert(actualSystem == "rocketmq", "期望messaging.system为 %s, 实际为 %s", "rocketmq", actualSystem)
	// 验证Span种类
	verifier.Assert(receive.SpanKind == trace.SpanKindConsumer, "期望为消费者Span, 实际为 %d", receive.SpanKind)

	verifier.Assert(process.Links[0].SpanContext.TraceID() == producter.SpanContext.TraceID(), "期望traceid为 %s, 实际为 %s", process.Links[0].SpanContext.TraceID(), producter.SpanContext.TraceID())
}
