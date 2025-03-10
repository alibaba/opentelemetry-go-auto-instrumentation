package ai

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"testing"
)

type testRequest struct {
	System    string
	Operation string
}
type testResponse struct {
}
type commonRequest struct{}

func (commonRequest) GetAIOperationName(request testRequest) string {
	return request.Operation
}
func (commonRequest) GetAISystem(request testRequest) string {
	return request.System
}

type ollamaRequest struct{}

func (ollamaRequest) GetAIRequestModel(request testRequest) string {
	return "deepseek:17b"
}
func (ollamaRequest) GetAIRequestEncodingFormats(request testRequest) []string {
	return []string{"string"}
}
func (ollamaRequest) GetAIRequestFrequencyPenalty(request testRequest) float64 {
	return 1.0
}
func (ollamaRequest) GetAIRequestPresencePenalty(request testRequest) float64 {
	return 1.0
}
func (ollamaRequest) GetAIResponseFinishReasons(request testRequest, response testResponse) []string {
	return []string{"stop"}
}
func (ollamaRequest) GetAIResponseModel(request testRequest, response testResponse) string {
	return "deepseek:17b"
}
func (ollamaRequest) GetAIRequestMaxTokens(request testRequest) int64 {
	return 10
}
func (ollamaRequest) GetAIUsageInputTokens(request testRequest) int64 {
	return 10
}
func (ollamaRequest) GetAIUsageOutputTokens(request testRequest, response testResponse) int64 {
	return 10
}
func (ollamaRequest) GetAIRequestStopSequences(request testRequest) []string {
	return []string{"stop"}
}
func (ollamaRequest) GetAIRequestTemperature(request testRequest) float64 {
	return 1.0
}
func (ollamaRequest) GetAIRequestTopK(request testRequest) float64 {
	return 1.0
}
func (ollamaRequest) GetAIRequestTopP(request testRequest) float64 {
	return 1.0
}
func (ollamaRequest) GetAIResponseID(request testRequest, response testResponse) string {
	return "chatcmpl-123"
}
func (ollamaRequest) GetAIServerAddress(request testRequest) string {
	return "127.0.0.1:1234"
}
func (ollamaRequest) GetAIRequestSeed(request testRequest) int64 {
	return 100
}

func TestCommonExtractorStart(t *testing.T) {
	Extractor := AICommonAttrsExtractor[testRequest, any, commonRequest]{}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs, _ = Extractor.OnStart(attrs, parentContext, testRequest{Operation: "llm", System: "langchain"})
	if len(attrs) == 0 {
		t.Fatal("attrs is empty")
	}
	if attrs[0].Key != "gen_ai.operation.name" || attrs[0].Value.AsString() != "llm" { //semconv.GenAIOperationNameKey
		t.Fatalf("gen_ai.operation.name be llm")
	}
	if attrs[1].Key != "gen_ai.system" || attrs[1].Value.AsString() != "langchain" { //semconv.GenAISystemKey
		t.Fatalf("gen_ai.system should be langchain")
	}
}
func TestCommonExtractorEnd(t *testing.T) {
	dbExtractor := AICommonAttrsExtractor[testRequest, any, commonRequest]{}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs, _ = dbExtractor.OnEnd(attrs, parentContext, testRequest{Operation: "llm", System: "langchain"}, testResponse{}, errors.New("test-err"))
	assert.Equal(t, semconv.ErrorTypeKey, attrs[0].Key)
	assert.Equal(t, "test-err", attrs[0].Value.AsString())
}

func TestAILLMAttrsExtractorStart(t *testing.T) {
	LLMExtractor := AILLMAttrsExtractor[testRequest, testResponse, commonRequest, ollamaRequest]{
		Base:      AICommonAttrsExtractor[testRequest, testResponse, commonRequest]{},
		LLMGetter: ollamaRequest{},
	}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs, _ = LLMExtractor.OnStart(attrs, parentContext, testRequest{Operation: "llm", System: "langchain"})
	if len(attrs) == 0 {
		t.Fatal("attrs is empty")
	}
	assert.Equal(t, "gen_ai.operation.name", attrs[0].Key) //semconv.GenAIOperationNameKey
	assert.Equal(t, "llm", attrs[0].Value.AsString())
	assert.Equal(t, "gen_ai.system", attrs[1].Key) //semconv.GenAISystemKey
	assert.Equal(t, "langchain", attrs[1].Value.AsString())

	assert.Equal(t, "gen_ai.request.model", attrs[2].Key) //semconv.GenAIRequestModelKey
	assert.Equal(t, "deepseek:17b", attrs[2].Value.AsString())
	assert.Equal(t, "gen_ai.request.encoding_formats", attrs[3].Key) //semconv.GenAIRequestEncodingFormatsKey
	assert.Equal(t, []string{"string"}, attrs[3].Value.AsStringSlice())
	assert.Equal(t, "gen_ai.request.max_tokens", attrs[4].Key) //semconv.GenAIRequestMaxTokensKey
	assert.Equal(t, int64(10), attrs[4].Value.AsInt64())
	assert.Equal(t, "gen_ai.request.frequency_penalty", attrs[5].Key) //semconv.GenAIRequestFrequencyPenaltyKey
	assert.Equal(t, 1.0, attrs[5].Value.AsFloat64())
	assert.Equal(t, "gen_ai.request.presence_penalty", attrs[6].Key) //semconv.GenAIRequestPresencePenaltyKey
	assert.Equal(t, 1.0, attrs[6].Value.AsFloat64())
	assert.Equal(t, "gen_ai.request.stop_sequences", attrs[7].Key) //semconv.GenAIRequestStopSequencesKey
	assert.Equal(t, []string{"stop"}, attrs[7].Value.AsStringSlice())
	assert.Equal(t, "gen_ai.request.temperature", attrs[8].Key) //semconv.GenAIRequestTemperatureKey
	assert.Equal(t, 1.0, attrs[8].Value.AsFloat64())
	assert.Equal(t, "gen_ai.request.top_k", attrs[9].Key) //semconv.GenAIRequestTopKKey
	assert.Equal(t, 1.0, attrs[9].Value.AsFloat64())
	assert.Equal(t, "gen_ai.request.top_p", attrs[10].Key) //semconv.GenAIRequestTopPKey
	assert.Equal(t, 1.0, attrs[10].Value.AsFloat64())
	assert.Equal(t, "gen_ai.usage.input_tokens", attrs[11].Key) //semconv.GenAIUsageInputTokensKey
	assert.Equal(t, int64(10), attrs[11].Value.AsInt64())
	assert.Equal(t, semconv.ServerAddressKey, attrs[12].Key)
	assert.Equal(t, "127.0.0.1:1234", attrs[12].Value.AsString())
	assert.Equal(t, "gen_ai.request.seed", attrs[13].Key) //semconv.GenAIRequestSeedKey
	assert.Equal(t, int64(100), attrs[13].Value.AsInt64())
}
func TestAILLMAttrsExtractorEnd(t *testing.T) {
	LLMExtractor := AILLMAttrsExtractor[testRequest, testResponse, commonRequest, ollamaRequest]{
		Base:      AICommonAttrsExtractor[testRequest, testResponse, commonRequest]{},
		LLMGetter: ollamaRequest{},
	}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs, _ = LLMExtractor.OnEnd(attrs, parentContext, testRequest{Operation: "llm", System: "langchain"}, testResponse{}, nil)
	assert.Equal(t, "gen_ai.response.finish_reasons", attrs[0].Key) //semconv.GenAIResponseFinishReasonsKey
	assert.Equal(t, []string{"stop"}, attrs[0].Value.AsStringSlice())
	assert.Equal(t, "gen_ai.response.id", attrs[1].Key) //semconv.GenAIResponseIDKey
	assert.Equal(t, "chatcmpl-123", attrs[1].Value.AsString())
	assert.Equal(t, "gen_ai.response.model", attrs[2].Key) //semconv.GenAIResponseModelKey
	assert.Equal(t, "deepseek:17b", attrs[2].Value.AsString())
	assert.Equal(t, "gen_ai.usage.output_tokens", attrs[3].Key) //semconv.GenAIUsageOutputTokensKey
	assert.Equal(t, int64(10), attrs[3].Value.AsInt64())
	assert.Equal(t, "gen_ai.response.id", attrs[4].Key) //semconv.GenAIResponseIDKey
	assert.Equal(t, "chatcmpl-123", attrs[4].Value.AsString())
}
