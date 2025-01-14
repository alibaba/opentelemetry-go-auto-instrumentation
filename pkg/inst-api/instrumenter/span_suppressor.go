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

	// database
	utils.DATABASE_SQL_SCOPE_NAME: utils.DB_CLIENT_KEY,
	utils.GO_REDIS_V9_SCOPE_NAME:  utils.DB_CLIENT_KEY,
	utils.GO_REDIS_V8_SCOPE_NAME:  utils.DB_CLIENT_KEY,
	utils.REDIGO_SCOPE_NAME:       utils.DB_CLIENT_KEY,
	utils.MONGO_SCOPE_NAME:        utils.DB_CLIENT_KEY,
	utils.GORM_SCOPE_NAME:         utils.DB_CLIENT_KEY,
}

var kindKey = map[string]trace.SpanKind{
	// http
	utils.FAST_HTTP_CLIENT_SCOPE_NAME:  trace.SpanKindClient,
	utils.FAST_HTTP_SERVER_SCOPE_NAME:  trace.SpanKindServer,
	utils.NET_HTTP_CLIENT_SCOPE_NAME:   trace.SpanKindClient,
	utils.NET_HTTP_SERVER_SCOPE_NAME:   trace.SpanKindServer,
	utils.HERTZ_HTTP_CLIENT_SCOPE_NAME: trace.SpanKindClient,
	utils.HERTZ_HTTP_SERVER_SCOPE_NAME: trace.SpanKindServer,

	// grpc
	utils.GRPC_CLIENT_SCOPE_NAME: trace.SpanKindClient,
	utils.GRPC_SERVER_SCOPE_NAME: trace.SpanKindServer,
	// kitex
	utils.KITEX_CLIENT_SCOPE_NAME: trace.SpanKindClient,
	utils.KITEX_SERVER_SCOPE_NAME: trace.SpanKindServer,

	// database
	utils.DATABASE_SQL_SCOPE_NAME: trace.SpanKindClient,
	utils.GO_REDIS_V9_SCOPE_NAME:  trace.SpanKindClient,
	utils.GO_REDIS_V8_SCOPE_NAME:  trace.SpanKindClient,
	utils.REDIGO_SCOPE_NAME:       trace.SpanKindClient,
	utils.MONGO_SCOPE_NAME:        trace.SpanKindClient,
	utils.GORM_SCOPE_NAME:         trace.SpanKindClient,
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
	return &SpanKindSuppressor{}
}

func (s *SpanKindSuppressor) StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	// do nothing
	return context
}

func (s *SpanKindSuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	span := trace.SpanFromContext(parentContext)
	if s, ok := span.(ottrace.ReadOnlySpan); ok {
		instScopeName := s.InstrumentationScope().Name
		if instScopeName != "" {
			parentSpanKind := kindKey[instScopeName]
			if spanKind != parentSpanKind {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
