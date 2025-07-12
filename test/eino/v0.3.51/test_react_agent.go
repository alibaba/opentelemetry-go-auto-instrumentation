package main

import (
	"context"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()
	g := compose.NewGraph[[]*schema.Message, *schema.Message]()
	reactAgentKeyOfLambda, err := NewMockReActAgentLambda(ctx)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("model", reactAgentKeyOfLambda)
	if err != nil {
		panic(err)
	}
	_ = g.AddEdge(compose.START, "model")
	_ = g.AddEdge("model", compose.END)
	graph, err := g.Compile(ctx)
	if err != nil {
		panic(err)
	}
	_, err = graph.Invoke(ctx, []*schema.Message{schema.UserMessage("hello")})
	if err != nil {
		panic(err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMAttributes(stubs[0][3], "chat", "eino", "mock-chat")
		verifier.VerifyLLMAttributes(stubs[0][6], "chat", "eino", "mock-chat")
		verifier.VerifyLLMCommonAttributes(stubs[0][10], "execute_tool", "eino", trace.SpanKindClient)
	}, 1)
}
