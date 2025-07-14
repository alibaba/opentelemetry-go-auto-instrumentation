// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package eino

import (
	"context"
	"reflect"
	_ "unsafe"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	utilscallbacks "github.com/cloudwego/eino/utils/callbacks"
)

//go:linkname newGraphOnEnter github.com/cloudwego/eino/compose.newGraphOnEnter
func newGraphOnEnter(call api.CallContext, cfg interface{}) {
	handler := utilscallbacks.NewHandlerHelper().
		Graph(NewComposeHandler("graph")).Chain(NewComposeHandler("chain")).Lambda(NewComposeHandler("lambda")).
		Prompt(einoPromptCallbackHandler()).Transformer(einoTransformCallbackHandler()).
		Embedding(einoEmbeddingCallbackHandler()).Indexer(einoIndexerCallbackHandler()).
		Retriever(einoRetrieverCallbackHandler()).Loader(einoLoaderCallbackHandler()).
		Tool(einoToolCallbackHandler()).ToolsNode(einoToolsNodeCallbackHandler()).
		Handler()
	callbacks.AppendGlobalHandlers(handler)
}

//go:linkname openaiGenerateOnEnter github.com/cloudwego/eino-ext/components/model/openai.openaiGenerateOnEnter
func openaiGenerateOnEnter(call api.CallContext, cm *openai.ChatModel, ctx context.Context, in []*schema.Message, opts ...model.Option) {
	cli := reflect.ValueOf(*cm).FieldByName("cli")
	config := ChatModelConfig{}
	if cli.IsValid() && !cli.IsNil() {
		conf := cli.Elem().FieldByName("config")
		if conf.IsValid() && !conf.IsNil() {
			if conf.Elem().FieldByName("BaseURL").IsValid() {
				config.BaseURL = conf.Elem().FieldByName("BaseURL").String()
			}
			if conf.Elem().FieldByName("PresencePenalty").IsValid() && !conf.Elem().FieldByName("PresencePenalty").IsNil() {
				config.PresencePenalty = conf.Elem().FieldByName("PresencePenalty").Elem().Float()
			}
			if conf.Elem().FieldByName("Seed").IsValid() && !conf.Elem().FieldByName("Seed").IsNil() {
				config.Seed = conf.Elem().FieldByName("Seed").Elem().Int()
			}
			if conf.Elem().FieldByName("FrequencyPenalty").IsValid() && !conf.Elem().FieldByName("FrequencyPenalty").IsNil() {
				config.FrequencyPenalty = conf.Elem().FieldByName("FrequencyPenalty").Elem().Float()
			}
		}
	}
	handler := utilscallbacks.NewHandlerHelper().ChatModel(einoModelCallHandler(config)).Handler()
	info := &callbacks.RunInfo{
		Name:      "OpenAI Stream",
		Type:      "OpenAI",
		Component: components.ComponentOfChatModel,
	}
	ctx = callbacks.InitCallbacks(ctx, info, handler)
	call.SetParam(1, ctx)
}

//go:linkname openaiStreamOnEnter github.com/cloudwego/eino-ext/components/model/openai.openaiStreamOnEnter
func openaiStreamOnEnter(call api.CallContext, cm *openai.ChatModel, ctx context.Context, in []*schema.Message, opts ...model.Option) {
	cli := reflect.ValueOf(*cm).FieldByName("cli")
	config := ChatModelConfig{}
	if cli.IsValid() && !cli.IsNil() {
		conf := cli.Elem().FieldByName("config")
		if conf.IsValid() && !conf.IsNil() {
			if conf.Elem().FieldByName("BaseURL").IsValid() {
				config.BaseURL = conf.Elem().FieldByName("BaseURL").String()
			}
			if conf.Elem().FieldByName("PresencePenalty").IsValid() && !conf.Elem().FieldByName("PresencePenalty").IsNil() {
				config.PresencePenalty = conf.Elem().FieldByName("PresencePenalty").Elem().Float()
			}
			if conf.Elem().FieldByName("Seed").IsValid() && !conf.Elem().FieldByName("Seed").IsNil() {
				config.Seed = conf.Elem().FieldByName("Seed").Elem().Int()
			}
			if conf.Elem().FieldByName("FrequencyPenalty").IsValid() && !conf.Elem().FieldByName("FrequencyPenalty").IsNil() {
				config.FrequencyPenalty = conf.Elem().FieldByName("FrequencyPenalty").Elem().Float()
			}
		}
	}
	handler := utilscallbacks.NewHandlerHelper().ChatModel(einoModelCallHandler(config)).Handler()
	info := &callbacks.RunInfo{
		Name:      "OpenAI Stream",
		Type:      "OpenAI",
		Component: components.ComponentOfChatModel,
	}
	ctx = callbacks.InitCallbacks(ctx, info, handler)
	call.SetParam(1, ctx)
}

//go:linkname ollamaGenerateOnEnter github.com/cloudwego/eino-ext/components/model/ollama.ollamaGenerateOnEnter
func ollamaGenerateOnEnter(call api.CallContext, cm *ollama.ChatModel, ctx context.Context, input []*schema.Message, opts ...model.Option) {
	conf := reflect.ValueOf(*cm).FieldByName("config")
	config := ChatModelConfig{}
	if conf.IsValid() && !conf.IsNil() {
		if conf.Elem().FieldByName("BaseURL").IsValid() {
			config.BaseURL = conf.Elem().FieldByName("BaseURL").String()
		}
	}
	handler := utilscallbacks.NewHandlerHelper().ChatModel(einoModelCallHandler(config)).Handler()
	info := &callbacks.RunInfo{
		Name:      "Ollama Generate",
		Type:      "Ollama",
		Component: components.ComponentOfChatModel,
	}
	ctx = callbacks.InitCallbacks(ctx, info, handler)
	call.SetParam(1, ctx)
}

//go:linkname ollamaStreamOnEnter github.com/cloudwego/eino-ext/components/model/ollama.ollamaStreamOnEnter
func ollamaStreamOnEnter(call api.CallContext, cm *ollama.ChatModel, ctx context.Context, input []*schema.Message, opts ...model.Option) {
	conf := reflect.ValueOf(*cm).FieldByName("config")
	config := ChatModelConfig{}
	if conf.IsValid() && !conf.IsNil() {
		if conf.Elem().FieldByName("BaseURL").IsValid() {
			config.BaseURL = conf.Elem().FieldByName("BaseURL").String()
		}
	}
	handler := utilscallbacks.NewHandlerHelper().ChatModel(einoModelCallHandler(config)).Handler()
	info := &callbacks.RunInfo{
		Name:      "Ollama Stream",
		Type:      "Ollama",
		Component: components.ComponentOfChatModel,
	}
	ctx = callbacks.InitCallbacks(ctx, info, handler)
	call.SetParam(1, ctx)
}
