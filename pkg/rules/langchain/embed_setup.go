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
	"github.com/tmc/langchaingo/embeddings"
)

// EmbedQuery
func singleEmbedOnEnter(call api.CallContext,
	e *embeddings.EmbedderImpl,
	ctx context.Context,
	text string,
) {
	request := langChainRequest{
		moduleName: MEmbedSingle,
		input: map[string]interface{}{
			"text": text,
		},
	}
	langCtx := langChainInstrument.Start(ctx, request)
	data := make(map[string]interface{})
	data["ctx"] = langCtx
	call.SetData(data)
}
func singleEmbedOnExit(
	call api.CallContext,
	emb []float32,
	err error,
) {
	request := langChainRequest{
		moduleName: MEmbedSingle,
	}
	data := call.GetData().(map[string]interface{})
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	if err != nil {
		langChainInstrument.End(ctx, request, nil, err)
		return
	}
	langChainInstrument.End(ctx, request, nil, nil)
}

// BatchedEmbed
func batchedEmbedOnEnter(call api.CallContext,
	ctx context.Context,
	embedder embeddings.EmbedderClient,
	texts []string,
	batchSize int,
) {
	request := langChainRequest{
		moduleName: MEmbedBatch,
		input: map[string]interface{}{
			"batchSize": batchSize,
		},
	}
	langCtx := langChainInstrument.Start(ctx, request)
	data := make(map[string]interface{})
	data["ctx"] = langCtx
	call.SetData(data)
}
func batchedEmbedOnExit(
	call api.CallContext,
	emb [][]float32,
	err error,
) {
	request := langChainRequest{
		moduleName: MEmbedBatch,
	}
	data := call.GetData().(map[string]interface{})
	ctx, ok := data["ctx"].(context.Context)
	if !ok {
		return
	}
	if err != nil {
		langChainInstrument.End(ctx, request, nil, err)
		return
	}
	langChainInstrument.End(ctx, request, nil, nil)
}
