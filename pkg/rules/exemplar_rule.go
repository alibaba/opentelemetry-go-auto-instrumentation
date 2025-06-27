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

package rules

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

// ExemplarContextRule captures trace context for exemplar association
type ExemplarContextRule struct {
	api.SkipRuleBase
}

func NewExemplarContextRule() *ExemplarContextRule {
	return &ExemplarContextRule{}
}

func (r *ExemplarContextRule) ID() string {
	return "exemplar_context_capture"
}

func (r *ExemplarContextRule) Version() string {
	return "v0.1.0"
}

func (r *ExemplarContextRule) Filter(call *api.CallContext) bool {
	// Apply to functions that:
	// 1. Have context.Context as first parameter
	// 2. Are HTTP handlers or gRPC handlers
	// 3. Start spans
	if len(call.Params) == 0 {
		return false
	}

	firstParam := call.Params[0]
	if firstParam.TypeName != "context.Context" {
		return false
	}

	// Check if it's an HTTP handler
	if call.HasParam("*http.Request") && call.HasParam("http.ResponseWriter") {
		return true
	}

	// Check if it's a gRPC handler
	if call.FuncName.Contains("grpc") && call.HasParam("context.Context") {
		return true
	}

	// Check if function starts a span
	if call.FuncName.Contains("Start") && call.FuncName.Contains("Span") {
		return true
	}

	return false
}

func (r *ExemplarContextRule) Apply(call *api.CallContext) {
	call.OnEnter(func(ctx *api.CallContext) {
		ctx.Set("exemplar_capture", true)
	})
}