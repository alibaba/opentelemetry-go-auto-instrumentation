package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func NewMockReActAgentLambda(ctx context.Context) (lba *compose.Lambda, err error) {
	ins, err := NewMockReactAgent(ctx)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}

func NewMockReactAgent(ctx context.Context) (*react.Agent, error) {
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mockBodyWithToolCall := `{"id":"e849dd51-7887-4e96-ac76-01e0b6dfa8b6","object":"chat.completion","created":1752200031,"model":"mock-chat","choices":[{"index":0,"message":{"role":"assistant","tool_calls":[{"index":0,"id":"call_0_b1bbef1e-376f-475c-b569-3f8d6da5fa44","type":"function","function":{"name":"greet","arguments":"{\"name\": \"test user\"}"}}]},"finish_reason":"tool_calls","content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false},"jailbreak":{"filtered":false,"detected":false},"profanity":{"filtered":false,"detected":false}}}],"usage":{"prompt_tokens":124,"completion_tokens":19,"total_tokens":143,"prompt_tokens_details":{"audio_tokens":0,"cached_tokens":64},"completion_tokens_details":null},"system_fingerprint":"fp_8802369eaa_prod0623_fp8_kvcache"}`
		mockBodyFinal := `{"id":"5f511bcc-d863-47e5-97a0-732d2a962c0f","object":"chat.completion","created":1752200038,"model":"mock-chat","choices":[{"index":0,"message":{"role":"assistant","content":"Hello form mock tool"},"finish_reason":"stop","content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false},"jailbreak":{"filtered":false,"detected":false},"profanity":{"filtered":false,"detected":false}}}],"usage":{"prompt_tokens":1779,"completion_tokens":340,"total_tokens":2119,"prompt_tokens_details":{"audio_tokens":0,"cached_tokens":1728},"completion_tokens_details":null},"system_fingerprint":"fp_8802369eaa_prod0623_fp8_kvcache"}`
		if count == 0 {
			w.Write([]byte(mockBodyWithToolCall))
			count++
		} else if count == 1 {
			w.Write([]byte(mockBodyFinal))
			count++
		}
	}))

	config := &openai.ChatModelConfig{
		BaseURL: ts.URL,
		Model:   "mock-chat",
	}
	cm, err := openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}

	mockTool := &MockGreetTool{}

	return react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{mockTool},
		},
		MaxStep: 3,
	})
}

func NewMockOllamaChatModelForInvoke(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mockBody := `{"model":"mock-chat","created_at":"2025-07-12T16:44:51.135504+08:00","message":{"role":"AI","content":"Hello! How can I assist you today?"},"done":true}`
		w.Write([]byte((mockBody)))
	}))
	config := &ollama.ChatModelConfig{
		BaseURL: ts.URL,
		Model:   "mock-chat",
	}
	cm, err = ollama.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func NewMockOllamaChatModelForStream(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)

		chunks := []string{
			`{"model":"mock-chat","created_at":"2025-07-12T16:44:51.135504+08:00","message":{"role":"assistant","content":"Hello"},"done":false}`,
			`{"model":"mock-chat","created_at":"2025-07-12T16:44:51.135504+08:00","message":{"role":"assistant","content":"! How can"},"done":false}`,
			`{"model":"mock-chat","created_at":"2025-07-12T16:44:51.135504+08:00","message":{"role":"assistant","content":" I assist you today?"},"done":false}`,
			`{"model":"mock-chat","created_at":"2025-07-12T16:44:51.135504+08:00","message":{"role":"assistant","content":""},"done":true}`,
		}

		for range chunks {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(10 * time.Millisecond)
		}
	}))

	config := &ollama.ChatModelConfig{
		BaseURL: ts.URL,
		Model:   "mock-chat",
	}
	cm, err = ollama.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func NewMockOpenAIChatModelForInvoke(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mockBody := `{"id":"5f511bcc-d863-47e5-97a0-732d2a962c0f","object":"chat.completion","created":1752200038,"model":"mock-chat","choices":[{"index":0,"message":{"role":"assistant","content":"Hello! How can I assist you today?"},"finish_reason":"stop","content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false},"jailbreak":{"filtered":false,"detected":false},"profanity":{"filtered":false,"detected":false}}}],"usage":{"prompt_tokens":1779,"completion_tokens":340,"total_tokens":2119,"prompt_tokens_details":{"audio_tokens":0,"cached_tokens":1728},"completion_tokens_details":null},"system_fingerprint":"fp_8802369eaa_prod0623_fp8_kvcache"}`
		w.Write([]byte((mockBody)))
	}))
	config := &openai.ChatModelConfig{
		BaseURL: ts.URL,
		Model:   "mock-chat",
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func NewMockOpenAIChatModelForStream(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		chunks := []string{
			`data: {"id":"5f511bcc-d863-47e5-97a0-732d2a962c0f","object":"chat.completion.chunk","created":1752200038,"model":"mock-chat","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`,
			`data: {"id":"5f511bcc-d863-47e5-97a0-732d2a962c0f","object":"chat.completion.chunk","created":1752200038,"model":"mock-chat","choices":[{"index":0,"delta":{"content":"! How can"},"finish_reason":null}]}`,
			`data: {"id":"5f511bcc-d863-47e5-97a0-732d2a962c0f","object":"chat.completion.chunk","created":1752200038,"model":"mock-chat","choices":[{"index":0,"delta":{"content":" I assist you today?"},"finish_reason":null}]}`,
			`data: {"id":"5f511bcc-d863-47e5-97a0-732d2a962c0f","object":"chat.completion.chunk","created":1752200038,"model":"mock-chat","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
			`data: [DONE]`,
		}

		for range chunks {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(10 * time.Millisecond)
		}
	}))

	config := &openai.ChatModelConfig{
		BaseURL: ts.URL,
		Model:   "mock-chat",
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

type MockGreetTool struct{}

func (t *MockGreetTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "greet",
		Desc: "greet with name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": {
					Desc:     "user name who to greet",
					Required: true,
					Type:     schema.String,
				},
			}),
	}, nil
}

func (t *MockGreetTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return "Hello from mock tool!", nil
}

type MockLoader struct{}

func (m *MockLoader) Load(ctx context.Context, src document.Source, opts ...document.LoaderOption) ([]*schema.Document, error) {
	return []*schema.Document{
		{ID: "doc1", Content: "这是第一个文档的内容"},
		{ID: "doc2", Content: "这是第二个文档的内容"},
		{ID: "doc3", Content: "这是第三个文档的内容"},
	}, nil
}

type MockEmbedder struct{}

func (m *MockEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))

	for i, text := range texts {
		embedding := make([]float64, 3)
		embedding[0] = float64(len(text))
		embedding[1] = float64(strings.Count(text, "文档"))
		embedding[2] = float64(i + 1)
		embeddings[i] = embedding
	}

	return embeddings, nil
}

type MockIndexer struct {
	storage map[string]*schema.Document
}

func NewMockIndexer() *MockIndexer {
	return &MockIndexer{
		storage: make(map[string]*schema.Document),
	}
}

func (m *MockIndexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	ids := make([]string, len(docs))

	for i, doc := range docs {
		m.storage[doc.ID] = doc
		ids[i] = doc.ID
	}

	return ids, nil
}

type MockRetriever struct {
	storage map[string]*schema.Document
}

func NewMockRetriever(storage map[string]*schema.Document) *MockRetriever {
	return &MockRetriever{storage: storage}
}

func (m *MockRetriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	var results []*schema.Document

	for _, doc := range m.storage {
		if strings.Contains(doc.Content, query) || strings.Contains(query, doc.ID) {
			results = append(results, doc)
		}
	}

	return results, nil
}
