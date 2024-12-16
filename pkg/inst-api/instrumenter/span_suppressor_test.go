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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestNoopSpanSuppressor(t *testing.T) {
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
			utils.HTTP_CLIENT_KEY,
		},
	}
	builder := Builder[testRequest, testResponse]{}
	builder.Init().SetSpanNameExtractor(testNameExtractor{}).
		SetSpanKindExtractor(&AlwaysClientExtractor[testRequest]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:      utils.FAST_HTTP_CLIENT_SCOPE_NAME,
			Version:   "test",
			SchemaURL: "test",
		})
	instrumenter := builder.BuildInstrumenter()
	ctx := context.Background()
	traceProvider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(traceProvider)
	newCtx := instrumenter.Start(ctx, testRequest{})
	span := trace.SpanFromContext(newCtx)
	newCtx = s.StoreInContext(newCtx, trace.SpanKindClient, span)
	if !s.ShouldSuppress(newCtx, trace.SpanKindClient) {
		t.Errorf("should suppress span")
	}
}

func TestSpanKeySuppressorNotMatch(t *testing.T) {
	s := SpanKeySuppressor{
		spanKeys: []attribute.Key{
			utils.RPC_CLIENT_KEY,
		},
	}
	builder := Builder[testRequest, testResponse]{}
	builder.Init().SetSpanNameExtractor(testNameExtractor{}).
		SetSpanKindExtractor(&AlwaysClientExtractor[testRequest]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:      utils.FAST_HTTP_CLIENT_SCOPE_NAME,
			Version:   "test",
			SchemaURL: "test",
		})
	instrumenter := builder.BuildInstrumenter()
	ctx := context.Background()
	traceProvider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(traceProvider)
	newCtx := instrumenter.Start(ctx, testRequest{})
	span := trace.SpanFromContext(newCtx)
	newCtx = s.StoreInContext(newCtx, trace.SpanKindClient, span)
	if s.ShouldSuppress(newCtx, trace.SpanKindClient) {
		t.Errorf("should not suppress span with different span key")
	}
}
