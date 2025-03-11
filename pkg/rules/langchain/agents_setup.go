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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func doActionOnEnter(call api.CallContext,
	e *agents.Executor,
	ctx context.Context,
	steps []schema.AgentStep,
	nameToTool map[string]tools.Tool,
	action schema.AgentAction,
) {
	buf := new(bytes.Buffer)
	if json.NewEncoder(buf).Encode(steps) != nil {
		buf.Reset()
		fmt.Fprintf(buf, "%#v", steps)
	}
	request := langChainRequest{
		operationName: MAgentAction,
		system:        "langchain",
		input: map[string]interface{}{
			"tool":       action.Tool,
			"tool-id":    action.ToolID,
			"tool-input": action.ToolInput,
			"log":        action.Log,
			"AgentStep":  buf.String(),
		},
	}
	langCtx := langChainCommonInstrument.Start(ctx, request)
	data := make(map[string]interface{})
	data["ctx"] = langCtx
	call.SetData(data)
}
func doActionOnExit(call api.CallContext, steps []schema.AgentStep, err error) {
	data := call.GetData().(map[string]interface{})
	buf := new(bytes.Buffer)
	request := langChainRequest{
		operationName: MAgentAction,
		system:        "langchain",
	}
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	if err != nil {
		langChainCommonInstrument.End(ctx, request, nil, err)
		return
	}
	if json.NewEncoder(buf).Encode(steps) != nil {
		request.output = map[string]any{"AgentStep": fmt.Sprintf("%#v", steps)}
	} else {
		request.output = map[string]any{"AgentStep": buf.String()}
	}
	langChainCommonInstrument.End(ctx, request, nil, nil)
}
