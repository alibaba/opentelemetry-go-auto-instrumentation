package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	/*llm, err := ollama.New([]ollama.Option{
		ollama.WithModel("qwen2.5:0.5b"),
		ollama.WithServerURL("http://127.0.0.1:" + os.Getenv("OLLAMA_QWEN_PORT")),
	}...)
	if err != nil {
		panic(err)
	}*/
	calc := new(tools.Calculator)
	ag := agents.NewConversationalAgent(fakeAgentLlm{}, []tools.Tool{getAgeTool{}, calc})
	ex := agents.NewExecutor(ag, agents.WithMaxIterations(100))
	_, err := ex.Call(context.Background(), map[string]any{"input": "请问张三的年龄"})
	if err != nil {
		panic(err)
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMCommonAttributes(stubs[1][0], "agentAction", "langchain", trace.SpanKindClient)
	}, 6)
}

type getAgeTool struct {
}

var _ tools.Tool = getAgeTool{}

func (c getAgeTool) Description() string {
	return `这是一个根据人名回答年龄的工具，如果需要查询一个人的年龄可以用这个工具。输入：人名，返回：年龄"`
}
func (c getAgeTool) Name() string {
	return "getAge"
}
func (c getAgeTool) Call(ctx context.Context, input string) (string, error) {
	return "20", nil
}

var step int

type fakeAgentLlm struct {
}

var _ llms.Model = fakeAgentLlm{}

func (f fakeAgentLlm) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	defer func() { step++ }()
	res := new(llms.ContentResponse)
	choice := new(llms.ContentChoice)
	if step == 0 {
		choice.Content = "Do I need to use a tool? Yes\nAction: getAge\nAction Input: 张三"
	} else if step == 1 {
		choice.Content = "No\nAI: 20"
	} else {
		choice.Content = "No\nAI: this not in fake llm"
	}
	fmt.Println(choice)
	res.Choices = append(res.Choices, choice)
	return res, nil
}
func (f fakeAgentLlm) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, f, prompt, options...)
}
