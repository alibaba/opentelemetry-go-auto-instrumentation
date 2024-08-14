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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SpanSuppressor interface {
	StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context
	ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool
}

type NoopSpanSuppressor struct {
}

func NewNoopSpanSuppressor() *NoopSpanSuppressor {
	return &NoopSpanSuppressor{}
}

func (n *NoopSpanSuppressor) StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	return context
}

func (n *NoopSpanSuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	return false
}

type SpanKeySuppressor struct {
	spanKeys []attribute.Key
}

func NewSpanKeySuppressor(spanKeys []attribute.Key) *SpanKeySuppressor {
	return &SpanKeySuppressor{spanKeys: spanKeys}
}

func (s *SpanKeySuppressor) StoreInContext(ctx context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	for _, spanKey := range s.spanKeys {
		ctx = context.WithValue(ctx, spanKey, span)
	}
	return ctx
}

func (s *SpanKeySuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	for _, spanKey := range s.spanKeys {
		if parentContext.Value(spanKey) == nil {
			return false
		}
	}
	return true
}

// TODO: semconv span suppressor
