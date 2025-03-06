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
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// GetRelevantDocuments
func getRelevantDocumentsOnEnter(call api.CallContext,
	r vectorstores.Retriever,
	ctx context.Context,
	query string,
) {
	request := langChainRequest{
		moduleName: MRelevantDoc,
		input: map[string]interface{}{
			"query": query,
		},
	}
	langCtx := langChainInstrument.Start(ctx, request)
	data := make(map[string]interface{})
	data["ctx"] = langCtx
	call.SetData(data)
}
func getRelevantDocumentsOnExit(
	call api.CallContext,
	s []schema.Document, err error,
) {
	buf := new(bytes.Buffer)
	if json.NewEncoder(buf).Encode(s) != nil {
		buf.Reset()
		fmt.Fprintf(buf, "%#v", s)
	}
	data := call.GetData().(map[string]interface{})
	request := langChainRequest{
		moduleName: MRelevantDoc,
		output: map[string]interface{}{
			"schema-doc": buf.String(),
		},
	}
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
