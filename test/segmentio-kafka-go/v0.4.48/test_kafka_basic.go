package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	ctx := context.Background()

	// Initialize producer with cleanup
	producer := initProducer()
	defer producer.Close()

	// Initialize consumer with cleanup
	consumer := initConsumer()
	defer consumer.Close()

	// Send message
	messages := []kafka.Message{kafka.Message{
		Value: []byte("hello world1"),
	}, kafka.Message{
		Value: []byte("hello world2"),
	}}
	if err := producer.WriteMessages(ctx, messages...); err != nil {
		panic(err)
	}

	// Read message
	_, err := consumer.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = consumer.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}

	// Verify OpenTelemetry traces
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyMQPublishAttributes(stubs[0][0], "", "", "", "publish", topicName, "kafka")
		verifier.VerifyMQConsumeAttributes(stubs[0][1], "", "", "", "process", topicName, "kafka")
		verifier.VerifyMQConsumeAttributes(stubs[0][2], "", "", "", "process", topicName, "kafka")
	}, 1)
}
