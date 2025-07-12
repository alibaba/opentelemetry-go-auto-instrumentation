package eino

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/ai"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type einoLLMAttrsGetter struct{}

var _ ai.LLMAttrsGetter[einoLLMRequest, einoLLMResponse] = einoLLMAttrsGetter{}
var _ ai.CommonAttrsGetter[einoLLMRequest, any] = einoLLMAttrsGetter{}

func (e einoLLMAttrsGetter) GetAIOperationName(request einoLLMRequest) string {
	return request.operationName
}

func (e einoLLMAttrsGetter) GetAISystem(request einoLLMRequest) string {
	return "eino"
}

func (e einoLLMAttrsGetter) GetAIRequestModel(request einoLLMRequest) string {
	return request.modelName
}

func (e einoLLMAttrsGetter) GetAIRequestEncodingFormats(request einoLLMRequest) []string {
	return request.encodingFormats
}

func (e einoLLMAttrsGetter) GetAIRequestFrequencyPenalty(request einoLLMRequest) float64 {
	return request.frequencyPenalty
}

func (e einoLLMAttrsGetter) GetAIRequestPresencePenalty(request einoLLMRequest) float64 {
	return request.frequencyPenalty
}

func (e einoLLMAttrsGetter) GetAIResponseFinishReasons(request einoLLMRequest, response einoLLMResponse) []string {
	return response.responseFinishReasons
}

func (e einoLLMAttrsGetter) GetAIResponseModel(request einoLLMRequest, response einoLLMResponse) string {
	return response.responseModel
}

func (e einoLLMAttrsGetter) GetAIRequestMaxTokens(request einoLLMRequest) int64 {
	return request.maxTokens
}

func (e einoLLMAttrsGetter) GetAIUsageInputTokens(request einoLLMRequest) int64 {
	return request.usageInputTokens
}

func (e einoLLMAttrsGetter) GetAIUsageOutputTokens(request einoLLMRequest, response einoLLMResponse) int64 {
	return response.usageOutputTokens
}

func (e einoLLMAttrsGetter) GetAIRequestStopSequences(request einoLLMRequest) []string {
	return request.stopSequences
}

func (e einoLLMAttrsGetter) GetAIRequestTemperature(request einoLLMRequest) float64 {
	return request.temperature
}

func (e einoLLMAttrsGetter) GetAIRequestTopK(request einoLLMRequest) float64 {
	return request.topK
}

func (e einoLLMAttrsGetter) GetAIRequestTopP(request einoLLMRequest) float64 {
	return request.topP
}

func (e einoLLMAttrsGetter) GetAIResponseID(request einoLLMRequest, response einoLLMResponse) string {
	return response.responseID
}

func (e einoLLMAttrsGetter) GetAIServerAddress(request einoLLMRequest) string {
	return request.serverAddress
}

func (e einoLLMAttrsGetter) GetAIRequestSeed(request einoLLMRequest) int64 {
	return request.seed
}

func BuildEinoLLMInstrumenter() instrumenter.Instrumenter[einoLLMRequest, einoLLMResponse] {
	builder := instrumenter.Builder[einoLLMRequest, einoLLMResponse]{}
	return builder.Init().SetSpanNameExtractor(&ai.AISpanNameExtractor[einoLLMRequest, einoLLMResponse]{Getter: einoLLMAttrsGetter{}}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[einoLLMRequest]{}).
		AddAttributesExtractor(&ai.AILLMAttrsExtractor[einoLLMRequest, einoLLMResponse, einoLLMAttrsGetter, einoLLMAttrsGetter]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.EINO_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
