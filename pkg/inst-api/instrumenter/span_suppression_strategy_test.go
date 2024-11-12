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

package instrumenter

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"testing"
)

func TestNoneShouldNotSuppressAnything(t *testing.T) {
	strategy := &NoneStrategy{}
	suppressor := strategy.create(nil)
	kind := trace.SpanKindClient
	spanKey := utils.HTTP_CLIENT_KEY
	context.WithValue(context.Background(), spanKey, noop.Span{})
	assert.False(t, suppressor.ShouldSuppress(context.Background(), kind))
}

func TestNoneShouldNotStoreSpansInContext(t *testing.T) {
	strategy := &NoneStrategy{}
	suppressor := strategy.create(nil)
	kind := trace.SpanKindClient
	ctx := context.Background()
	ctxNew := suppressor.StoreInContext(ctx, kind, noop.Span{})
	assert.DeepEqual(t, ctx, ctxNew)
}

func TestSemConvShouldSuppressContextWhenAllSpanKeysArePresent(t *testing.T) {
	spanKeys := []attribute.Key{utils.DB_CLIENT_KEY, utils.RPC_CLIENT_KEY}
	strategy := &SemConvStrategy{}
	suppressor := strategy.create(spanKeys)
	ctx := context.WithValue(context.WithValue(context.Background(), utils.DB_CLIENT_KEY, noop.Span{}), utils.RPC_CLIENT_KEY, noop.Span{})
	assert.True(t, suppressor.ShouldSuppress(ctx, trace.SpanKindServer))
}

func TestSemConvShouldNotSuppressContextWithPartiallyDifferentSpanKeys(t *testing.T) {
	spanKeys := []attribute.Key{utils.DB_CLIENT_KEY, utils.RPC_CLIENT_KEY}
	strategy := &SemConvStrategy{}
	suppressor := strategy.create(spanKeys)
	ctx := context.WithValue(context.WithValue(context.Background(), utils.DB_CLIENT_KEY, noop.Span{}), utils.HTTP_CLIENT_KEY, noop.Span{})
	assert.False(t, suppressor.ShouldSuppress(ctx, trace.SpanKindServer))
}

func TestSpanKindShouldSuppressSameKind(t *testing.T) {
	strategy := &SpanKindStrategy{}
	suppressor := strategy.create(nil)
	root := context.Background()
	ctx := suppressor.StoreInContext(root, trace.SpanKindServer, noop.Span{})
	assert.True(t, suppressor.ShouldSuppress(ctx, trace.SpanKindServer))
}
