package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"time"
)

func main() {
	// 初始化生产者
	initTopic()
	p := initRocketMQ()
	defer p.Shutdown()

	// 初始化消费者
	c := initConsumer()

	// 发送消息
	msg := primitive.NewMessage(topicName, []byte("Hello RocketMQ"))
	msg.WithTag("test_tag")

	result, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	log.Printf("消息发送成功: %s\n", result.String())

	// 注册消息处理函数
	err = c.Subscribe(topicName, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			log.Printf("收到消息: %s\n", string(msg.Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}

	// 启动消费者
	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Shutdown()

	time.Sleep(10 * time.Second)

	// 验证OpenTelemetry跟踪
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		VerifyRocketMQAttributes(stubs[0][0], stubs[0][1], topicName, "test_tag")
	}, 1)
}
