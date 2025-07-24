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
	"context"
	_ "unsafe" // Required for go:linkname
	
	"github.com/alibaba/loongsuite-go-agent/pkg/api"
	ollamaapi "github.com/ollama/ollama/api"
)

// GENERATE API HOOKS

//go:linkname clientGenerateOnEnter github.com/ollama/ollama/api.clientGenerateOnEnter
func clientGenerateOnEnter(call api.CallContext, c *ollamaapi.Client, ctx context.Context, req *ollamaapi.GenerateRequest, fn ollamaapi.GenerateResponseFunc) {
	// Create request tracking object
	ollamaReq := ollamaRequest{
		operationType: "generate",
		model:        req.Model,
		prompt:       req.Prompt,
	}
	
	// Start OpenTelemetry span
	ctx = ollamaInstrumenter.Start(ctx, ollamaReq)
	
	// Update context parameter
	call.SetParam(1, ctx)
	
	// CRITICAL: Wrap the callback to capture response data
	var finalResponse ollamaapi.GenerateResponse
	var wrappedFn ollamaapi.GenerateResponseFunc = func(resp ollamaapi.GenerateResponse) error {
		// Always update with latest response
		// The final response will have Done=true
		if resp.Done {
			finalResponse = resp
		}
		
		// Call the original callback if provided
		if fn != nil {
			return fn(resp)
		}
		return nil
	}
	
	// Replace the callback parameter with our wrapped version
	call.SetParam(3, wrappedFn)
	
	// Store context and response pointer for OnExit
	data := make(map[string]interface{})
	data["ctx"] = ctx
	data["request"] = &ollamaReq
	data["finalResponsePtr"] = &finalResponse
	call.SetData(data)
}

//go:linkname clientGenerateOnExit github.com/ollama/ollama/api.clientGenerateOnExit
func clientGenerateOnExit(call api.CallContext, err error) {
	// Retrieve data stored in OnEnter
	data, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	
	// Get context from data
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	
	// Get request from data
	reqPtr, ok := data["request"].(*ollamaRequest)
	if !ok || reqPtr == nil {
		return
	}
	
	// Create response object
	ollamaResp := ollamaResponse{
		err: err,
	}
	
	// Extract response data if no error
	if err == nil {
		// Get the final response captured by our wrapped callback
		if respPtr, ok := data["finalResponsePtr"].(*ollamaapi.GenerateResponse); ok && respPtr != nil {
			ollamaResp.promptTokens = respPtr.PromptEvalCount
			ollamaResp.completionTokens = respPtr.EvalCount
			ollamaResp.content = respPtr.Response
			
			// Update request with token counts for the instrumenter
			reqPtr.promptTokens = respPtr.PromptEvalCount
			reqPtr.completionTokens = respPtr.EvalCount
		}
	}
	
	// End OpenTelemetry span
	ollamaInstrumenter.End(ctx, *reqPtr, ollamaResp, err)
}

// CHAT API HOOKS

//go:linkname clientChatOnEnter github.com/ollama/ollama/api.clientChatOnEnter
func clientChatOnEnter(call api.CallContext, c *ollamaapi.Client, ctx context.Context, req *ollamaapi.ChatRequest, fn ollamaapi.ChatResponseFunc) {
	// Create request tracking object
	ollamaReq := ollamaRequest{
		operationType: "chat",
		model:        req.Model,
		messages:     req.Messages,
	}
	
	// Start OpenTelemetry span
	ctx = ollamaInstrumenter.Start(ctx, ollamaReq)
	
	// Update context parameter
	call.SetParam(1, ctx)
	
	// CRITICAL: Wrap the callback to capture response data
	var finalResponse ollamaapi.ChatResponse
	var wrappedFn ollamaapi.ChatResponseFunc = func(resp ollamaapi.ChatResponse) error {
		// Always update with latest response
		// The final response will have Done=true
		if resp.Done {
			finalResponse = resp
		}
		
		// Call the original callback if provided
		if fn != nil {
			return fn(resp)
		}
		return nil
	}
	
	// Replace the callback parameter with our wrapped version
	call.SetParam(3, wrappedFn)
	
	// Store context and response pointer for OnExit
	data := make(map[string]interface{})
	data["ctx"] = ctx
	data["request"] = &ollamaReq
	data["finalResponsePtr"] = &finalResponse
	call.SetData(data)
}

//go:linkname clientChatOnExit github.com/ollama/ollama/api.clientChatOnExit
func clientChatOnExit(call api.CallContext, err error) {
	// Retrieve data stored in OnEnter
	data, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	
	// Get context from data
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	
	// Get request from data
	reqPtr, ok := data["request"].(*ollamaRequest)
	if !ok || reqPtr == nil {
		return
	}
	
	// Create response object
	ollamaResp := ollamaResponse{
		err: err,
	}
	
	// Extract response data if no error
	if err == nil {
		// Get the final response captured by our wrapped callback
		if respPtr, ok := data["finalResponsePtr"].(*ollamaapi.ChatResponse); ok && respPtr != nil {
			// Token counts are in embedded Metrics struct
			ollamaResp.promptTokens = respPtr.PromptEvalCount
			ollamaResp.completionTokens = respPtr.EvalCount
			ollamaResp.content = respPtr.Message.Content
			
			// Update request with token counts for the instrumenter
			reqPtr.promptTokens = respPtr.PromptEvalCount
			reqPtr.completionTokens = respPtr.EvalCount
		}
	}
	
	// End OpenTelemetry span
	ollamaInstrumenter.End(ctx, *reqPtr, ollamaResp, err)
}