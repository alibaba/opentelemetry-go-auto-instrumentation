// Copyright (c) 2025 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ollama

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api-semconv/instrumenter/ai"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/instrumenter"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/version"
)

const (
	OLLAMA_SCOPE_NAME = "github.com/alibaba/loongsuite-go-agent/pkg/rules/ollama"
)

// ollamaAttrsGetter implements the interfaces for extracting attributes
type ollamaAttrsGetter struct{}

// Request attribute extraction methods
func (o ollamaAttrsGetter) GetAISystem(request ollamaRequest) string {
	return "ollama"
}

func (o ollamaAttrsGetter) GetAIRequestModel(request ollamaRequest) string {
	return request.model
}

func (o ollamaAttrsGetter) GetAIRequestTemperature(request ollamaRequest) float64 {
	// Temperature parameter not captured in this implementation
	return 0
}

func (o ollamaAttrsGetter) GetAIRequestMaxTokens(request ollamaRequest) int64 {
	// Max tokens parameter not captured in this implementation
	return 0
}

func (o ollamaAttrsGetter) GetAIRequestTopP(request ollamaRequest) float64 {
	// TopP parameter not captured in this implementation
	return 0
}

func (o ollamaAttrsGetter) GetAIRequestTopK(request ollamaRequest) float64 {
	// TopK parameter not captured in this implementation
	return 0
}

func (o ollamaAttrsGetter) GetAIRequestStopSequences(request ollamaRequest) []string {
	// Stop sequences not captured in this implementation
	return nil
}

func (o ollamaAttrsGetter) GetAIRequestFrequencyPenalty(request ollamaRequest) float64 {
	// Frequency penalty parameter not captured in this implementation
	return 0
}

func (o ollamaAttrsGetter) GetAIRequestPresencePenalty(request ollamaRequest) float64 {
	// Presence penalty parameter not captured in this implementation
	return 0
}

func (o ollamaAttrsGetter) GetAIRequestIsStream(request ollamaRequest) bool {
	// Return true if this is a streaming request
	return request.isStreaming
}

func (o ollamaAttrsGetter) GetAIOperationName(request ollamaRequest) string {
	return request.operationType
}

func (o ollamaAttrsGetter) GetAIRequestEncodingFormats(request ollamaRequest) []string {
	// Encoding formats not captured in this implementation
	return nil
}

func (o ollamaAttrsGetter) GetAIRequestSeed(request ollamaRequest) int64 {
	// Seed parameter not captured in this implementation
	return 0
}

// Response attribute extraction methods
func (o ollamaAttrsGetter) GetAIResponseModel(request ollamaRequest, response ollamaResponse) string {
	// Model comes from request
	return request.model
}

func (o ollamaAttrsGetter) GetAIUsageInputTokens(request ollamaRequest) int64 {
	return int64(request.promptTokens)
}

func (o ollamaAttrsGetter) GetAIUsageOutputTokens(request ollamaRequest, response ollamaResponse) int64 {
	return int64(request.completionTokens)
}

func (o ollamaAttrsGetter) GetAIResponseFinishReasons(request ollamaRequest, response ollamaResponse) []string {
	if response.err != nil {
		return []string{"error"}
	}
	return []string{"stop"}
}

func (o ollamaAttrsGetter) GetAIResponseID(request ollamaRequest, response ollamaResponse) string {
	// Response ID not available in Ollama API
	return ""
}

func (o ollamaAttrsGetter) GetAIServerAddress(request ollamaRequest) string {
	// Server address not captured in this implementation
	return ""
}

// BuildOllamaLLMInstrumenter creates the instrumenter using the generic pattern
func BuildOllamaLLMInstrumenter() instrumenter.Instrumenter[ollamaRequest, ollamaResponse] {
	builder := instrumenter.Builder[ollamaRequest, ollamaResponse]{}
	getter := ollamaAttrsGetter{}
	
	return builder.Init().
		SetSpanNameExtractor(&ai.AISpanNameExtractor[ollamaRequest, ollamaResponse]{Getter: getter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[ollamaRequest]{}).
		AddAttributesExtractor(&ai.AILLMAttrsExtractor[ollamaRequest, ollamaResponse, ollamaAttrsGetter, ollamaAttrsGetter]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    OLLAMA_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}

// Singleton instance
var ollamaInstrumenter = BuildOllamaLLMInstrumenter()