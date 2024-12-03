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

	"go.opentelemetry.io/otel/attribute"
	ottrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var scopeKey = map[string]attribute.Key{
	// http
	utils.FAST_HTTP_CLIENT_SCOPE_NAME:  utils.HTTP_CLIENT_KEY,
	utils.FAST_HTTP_SERVER_SCOPE_NAME:  utils.HTTP_SERVER_KEY,
	utils.NET_HTTP_CLIENT_SCOPE_NAME:   utils.HTTP_CLIENT_KEY,
	utils.NET_HTTP_SERVER_SCOPE_NAME:   utils.HTTP_SERVER_KEY,
	utils.HERTZ_HTTP_CLIENT_SCOPE_NAME: utils.HTTP_CLIENT_KEY,
	utils.HERTZ_HTTP_SERVER_SCOPE_NAME: utils.HTTP_SERVER_KEY,

	// grpc
	utils.GRPC_CLIENT_SCOPE_NAME: utils.RPC_CLIENT_KEY,
	utils.GRPC_SERVER_SCOPE_NAME: utils.RPC_SERVER_KEY,

	//utils.KRATOS_GRPC_INTERNAL_SCOPE_NAME: utils.RPC_CLIENT_KEY,
	//utils.KRATOS_HTTP_INTERNAL_SCOPE_NAME: utils.HTTP_CLIENT_KEY,

	// database
	utils.DATABASE_SQL_SCOPE_NAME: utils.DB_CLIENT_KEY,
	utils.GO_REDIS_V9_SCOPE_NAME:  utils.DB_CLIENT_KEY,
	utils.GO_REDIS_V8_SCOPE_NAME:  utils.DB_CLIENT_KEY,
	utils.REDIGO_SCOPE_NAME:       utils.DB_CLIENT_KEY,
	utils.MONGO_SCOPE_NAME:        utils.DB_CLIENT_KEY,
	utils.GORM_SCOPE_NAME:         utils.DB_CLIENT_KEY,
}

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
	// do nothing
	return ctx
}

func (s *SpanKeySuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	for _, spanKey := range s.spanKeys {
		span := trace.SpanFromContext(parentContext)
		if s, ok := span.(ottrace.ReadOnlySpan); ok {
			instScopeName := s.InstrumentationScope().Name
			if instScopeName != "" {
				parentSpanKey := scopeKey[instScopeName]
				if spanKey != parentSpanKey {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func NewSpanKindSuppressor() *SpanKindSuppressor {
	var m = make(map[trace.SpanKind]SpanSuppressor)
	m[trace.SpanKindServer] = NewSpanKeySuppressor([]attribute.Key{utils.KIND_SERVER})
	m[trace.SpanKindClient] = NewSpanKeySuppressor([]attribute.Key{utils.KIND_CLIENT})
	m[trace.SpanKindProducer] = NewSpanKeySuppressor([]attribute.Key{utils.KIND_PRODUCER})
	m[trace.SpanKindConsumer] = NewSpanKeySuppressor([]attribute.Key{utils.KIND_CONSUMER})

	return &SpanKindSuppressor{
		delegates: m,
	}
}

func (s *SpanKindSuppressor) StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	spanSuppressor, exists := s.delegates[spanKind]
	if !exists {
		return context
	}
	return spanSuppressor.StoreInContext(context, spanKind, span)
}

func (s *SpanKindSuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	spanSuppressor, exists := s.delegates[spanKind]
	if !exists {
		return false
	}
	return spanSuppressor.ShouldSuppress(parentContext, spanKind)
}
