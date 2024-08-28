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
//go:build ignore

package rule

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type dubboRequest struct {
	method   string
	addr     string
	metadata map[string]interface{}
}

type dubboResponse struct {
	statusCode string
}

func GenerateSpanName(invoker protocol.Invoker, inv protocol.Invocation) string {
	group := invoker.GetURL().GetParam("group", "")
	return group + invoker.GetURL().Path + "/" + inv.MethodName()
}

type metadataSupplier struct {
	metadata map[string]interface{}
}

var _ propagation.TextMapCarrier = &metadataSupplier{}

func (s *metadataSupplier) Get(key string) string {
	if s.metadata == nil {
		return ""
	}
	item, ok := s.metadata[key].([]string)
	if !ok {
		return ""
	}
	if len(item) == 0 {
		return ""
	}
	return item[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	if s.metadata == nil {
		s.metadata = map[string]interface{}{}
	}
	s.metadata[key] = value
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(s.metadata))
	for key := range s.metadata {
		out = append(out, key)
	}
	return out
}

// Inject injects correlation context and span context into the dubbo
// metadata object. This function is meant to be used on outgoing
// requests.
func Inject(ctx context.Context, metadata map[string]interface{}, propagators propagation.TextMapPropagator) {
	propagators.Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract returns the baggage and span context that
// another service encoded in the dubbo metadata object with Inject.
// This function is meant to be used on incoming requests.
func Extract(ctx context.Context, metadata map[string]interface{}, propagators propagation.TextMapPropagator) (baggage.Baggage, trace.SpanContext) {
	ctx = propagators.Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})
	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}
