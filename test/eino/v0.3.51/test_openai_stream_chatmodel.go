package main

import (
	"context"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/cloudwego/eino/schema"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	ctx := context.Background()
	cm, _ := NewMockOpenAIChatModelForStream(ctx)
	_, err := cm.Stream(ctx, []*schema.Message{schema.UserMessage("Hello")})
	if err != nil {
		panic(err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMAttributes(stubs[0][0], "chat", "eino", "mock-chat")
	}, 1)
}
