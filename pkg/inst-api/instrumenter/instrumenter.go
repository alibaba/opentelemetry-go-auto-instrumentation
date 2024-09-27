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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Instrumenter[REQUEST any, RESPONSE any] interface {
	StartAndEnd(parentContext context.Context, request REQUEST, response RESPONSE, err error, startTime, endTime time.Time)
	Start(parentContext context.Context, request REQUEST) context.Context
	End(ctx context.Context, request REQUEST, response RESPONSE, err error)
}

type InternalInstrumenter[REQUEST any, RESPONSE any] struct {
	enabler              InstrumentEnabler
	spanNameExtractor    SpanNameExtractor[REQUEST]
	spanKindExtractor    SpanKindExtractor[REQUEST]
	spanStatusExtractor  SpanStatusExtractor[REQUEST, RESPONSE]
	attributesExtractors []AttributesExtractor[REQUEST, RESPONSE]
	operationListeners   []*OperationListenerWrapper
	operationMetrics     []OperationMetrics
	contextCustomizers   []ContextCustomizer[REQUEST]
	spanSuppressor       SpanSuppressor
	tracer               trace.Tracer
	instVersion          string
}

type PropagatingToDownstreamInstrumenter[REQUEST any, RESPONSE any] struct {
	carrierGetter func(REQUEST) propagation.TextMapCarrier
	prop          propagation.TextMapPropagator
	base          InternalInstrumenter[REQUEST, RESPONSE]
}

type PropagatingFromUpstreamInstrumenter[REQUEST any, RESPONSE any] struct {
	carrierGetter func(REQUEST) propagation.TextMapCarrier
	prop          propagation.TextMapPropagator
	base          InternalInstrumenter[REQUEST, RESPONSE]
}

func (i *InternalInstrumenter[REQUEST, RESPONSE]) StartAndEnd(parentContext context.Context, request REQUEST, response RESPONSE, err error, startTime, endTime time.Time) {
	ctx := i.doStart(parentContext, request, startTime)
	i.doEnd(ctx, request, response, err, endTime)
}

func (i *InternalInstrumenter[REQUEST, RESPONSE]) Start(parentContext context.Context, request REQUEST) context.Context {
	return i.doStart(parentContext, request, time.Now())
}

func (i *InternalInstrumenter[REQUEST, RESPONSE]) doStart(parentContext context.Context, request REQUEST, timestamp time.Time) context.Context {
	if i.enabler != nil && !i.enabler.IsEnabled() {
		return parentContext
	}
	for _, listener := range i.operationListeners {
		parentContext = listener.OnBeforeStart(parentContext, timestamp)
	}
	// extract span name
	spanName := i.spanNameExtractor.Extract(request)
	spanKind := i.spanKindExtractor.Extract(request)
	newCtx, span := i.tracer.Start(parentContext, spanName, trace.WithSpanKind(spanKind))
	attrs := make([]attribute.KeyValue, 0, 20)
	// extract span attrs
	for _, extractor := range i.attributesExtractors {
		attrs = extractor.OnStart(attrs, newCtx, request)
	}
	// execute context customizer hook
	for _, customizer := range i.contextCustomizers {
		newCtx = customizer.OnStart(newCtx, request, attrs)
	}
	for _, listener := range i.operationListeners {
		newCtx = listener.OnBeforeEnd(newCtx, attrs, timestamp)
	}
	span.SetAttributes(attrs...)
	return i.spanSuppressor.StoreInContext(newCtx, spanKind, span)
}

func (i *InternalInstrumenter[REQUEST, RESPONSE]) End(ctx context.Context, request REQUEST, response RESPONSE, err error) {
	i.doEnd(ctx, request, response, err, time.Now())
}

func (i *InternalInstrumenter[REQUEST, RESPONSE]) doEnd(ctx context.Context, request REQUEST, response RESPONSE, err error, timestamp time.Time) {
	if i.enabler != nil && !i.enabler.IsEnabled() {
		return
	}
	for _, listener := range i.operationListeners {
		listener.OnAfterStart(ctx, timestamp)
	}
	span := trace.SpanFromContext(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	var attrs []attribute.KeyValue
	// extract span attributes
	for _, extractor := range i.attributesExtractors {
		attrs = extractor.OnEnd(attrs, ctx, request, response, err)
	}
	i.spanStatusExtractor.Extract(span, request, response, err)
	span.SetAttributes(attrs...)
	span.End(trace.WithTimestamp(timestamp))
	for _, listener := range i.operationListeners {
		listener.OnAfterEnd(ctx, attrs, timestamp)
	}
}

func (p *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]) StartAndEnd(parentContext context.Context, request REQUEST, response RESPONSE, err error, startTime, endTime time.Time) {
	newCtx := p.base.doStart(parentContext, request, startTime)
	if p.carrierGetter != nil {
		if p.prop != nil {
			p.prop.Inject(newCtx, p.carrierGetter(request))
		} else {
			otel.GetTextMapPropagator().Inject(newCtx, p.carrierGetter(request))
		}
	}
	p.base.doEnd(newCtx, request, response, err, endTime)
}

func (p *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]) Start(parentContext context.Context, request REQUEST) context.Context {
	newCtx := p.base.Start(parentContext, request)
	if p.carrierGetter != nil {
		if p.prop != nil {
			p.prop.Inject(newCtx, p.carrierGetter(request))
		} else {
			otel.GetTextMapPropagator().Inject(newCtx, p.carrierGetter(request))
		}

	}
	return newCtx
}

func (p *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]) End(ctx context.Context, request REQUEST, response RESPONSE, err error) {
	p.base.End(ctx, request, response, err)
}

func (p *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]) StartAndEnd(parentContext context.Context, request REQUEST, response RESPONSE, err error, startTime, endTime time.Time) {
	var ctx context.Context
	if p.carrierGetter != nil {
		var extracted context.Context
		if p.prop != nil {
			extracted = p.prop.Extract(parentContext, p.carrierGetter(request))
		} else {
			extracted = otel.GetTextMapPropagator().Extract(parentContext, p.carrierGetter(request))
		}
		ctx = p.base.doStart(extracted, request, startTime)
	} else {
		ctx = parentContext
	}
	p.base.doEnd(ctx, request, response, err, endTime)
}

func (p *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]) Start(parentContext context.Context, request REQUEST) context.Context {
	if p.carrierGetter != nil {
		var extracted context.Context
		if p.prop != nil {
			extracted = p.prop.Extract(parentContext, p.carrierGetter(request))
		} else {
			extracted = otel.GetTextMapPropagator().Extract(parentContext, p.carrierGetter(request))
		}
		return p.base.Start(extracted, request)
	} else {
		return parentContext
	}
}

func (p *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]) End(ctx context.Context, request REQUEST, response RESPONSE, err error) {
	p.base.End(ctx, request, response, err)
}
