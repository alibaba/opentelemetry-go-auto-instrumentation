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
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	_ "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"time"
)

// 自定义不同的消费组常量
const (
	msgCount = 2
)

func main() {
	// 初始化生产者
	initTopic()
	p := initRocketMQ()
	defer p.Shutdown()

	// 测试集群消费模式
	clusterConsumer := initConsumer()

	err := clusterConsumer.Subscribe(topicName, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				log.Printf("集群模式消费: %s, Tags: %s\n", string(msg.Body), msg.GetTags())
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		panic(fmt.Errorf("订阅主题失败(集群模式)：%v", err))
	}

	// 等待消费者准备就绪
	time.Sleep(2 * time.Second)

	// 批量发送 处理模式为receive
	msgs := make([]*primitive.Message, msgCount)
	for i := 0; i < msgCount; i++ {
		msg := &primitive.Message{
			Topic: topicName,
			Body:  []byte(fmt.Sprintf("消费模式测试消息 %d", i)),
		}
		msg.WithTag(fmt.Sprintf("Tag%d", i))
		msgs[i] = msg
	}
	// 同步发送消息
	result, err := p.SendSync(context.Background(), msgs...)
	if err != nil {
		panic(fmt.Errorf("发送消息失败：%v", err))
	}
	fmt.Printf("消息发送成功: %s\n", result.MsgID)

	// 启动消费者
	err = clusterConsumer.Start()
	if err != nil {
		panic(fmt.Errorf("启动消费者失败(集群模式)：%v", err))
	}
	defer clusterConsumer.Shutdown()

	// 等待消费
	time.Sleep(10 * time.Second)

	//验证OpenTelemetry跟踪
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		VerifyRocketMQReceive(stubs[0][0], stubs[1][0], stubs[1][1])
		VerifyRocketMQConsumeAttributes(stubs[1][1], topicName, "Tag0", "", "process", false)
		VerifyRocketMQConsumeAttributes(stubs[1][2], topicName, "Tag1", "", "process", false)
	}, 1)
}
