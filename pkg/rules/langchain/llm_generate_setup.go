// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package langchain

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/tmc/langchaingo/llms"
)

func generateFromSinglePromptOnEnter(call api.CallContext,
	ctx context.Context,
	llm llms.Model,
	prompt string,
	options ...llms.CallOption,
) {
	request := langChainRequest{
		moduleName: MLlmGenerateSingle,
		system:     "langchain",
		input: map[string]any{
			"prompt": prompt,
		},
	}
	langCtx := langChainCommonInstrument.Start(ctx, request)
	data := make(map[string]interface{})
	data["ctx"] = langCtx
	call.SetData(data)
}
func generateFromSinglePromptOnExit(call api.CallContext, v string, err error) {
	data := call.GetData().(map[string]interface{})
	request := langChainRequest{
		moduleName: MLlmGenerateSingle,
		system:     "langchain",
		output:     map[string]any{"output": v},
	}
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	if err != nil {
		langChainCommonInstrument.End(ctx, request, nil, err)
		return
	}
	langChainCommonInstrument.End(ctx, request, nil, nil)
}
