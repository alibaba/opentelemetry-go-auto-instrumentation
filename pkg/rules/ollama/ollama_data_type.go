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
	"github.com/ollama/ollama/api"
)

// ollamaRequest represents an Ollama API request
type ollamaRequest struct {
	operationType string        // "chat" or "generate"
	model        string         // Model name
	messages     []api.Message  // For chat requests
	prompt       string         // For generate requests
	
	// Token counts - populated from response
	promptTokens     int
	completionTokens int
}

// ollamaResponse represents an Ollama API response
type ollamaResponse struct {
	// Token counts from the response
	promptTokens     int
	completionTokens int
	
	// Response content
	content string
	
	// Error if any
	err error
}