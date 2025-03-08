package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/fake"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	llm := fake.NewFakeLLM([]string{"你好"})
	_, err := llms.GenerateFromSinglePrompt(context.Background(), llm, "30的3次方是多少")
	if err != nil {
		panic(err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMAttributes(stubs[0][0], "llmGenerateSingle")
	}, 1)
}
