// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package instrumenter

import (
	"context"
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestNoopSpanSupressor(t *testing.T) {
	n := NoopSpanSuppressor{}
	ctx := context.Background()
	n.StoreInContext(ctx, trace.SpanKindClient, noop.Span{})
	if n.ShouldSuppress(ctx, trace.SpanKindClient) != false {
		t.Errorf("should not suppress span")
	}
}

func TestSpanKeySuppressor(t *testing.T) {
	s := SpanKeySuppressor{
		spanKeys: []attribute.Key{
			utils.RPC_CLIENT_KEY,
		},
	}
	ctx := context.Background()
	newCtx := s.StoreInContext(ctx, trace.SpanKindClient, noop.Span{})
	if !s.ShouldSuppress(newCtx, trace.SpanKindClient) {
		t.Errorf("should suppress span")
	}
	if s.ShouldSuppress(context.Background(), trace.SpanKindClient) {
		t.Errorf("should not suppress span")
	}
}
