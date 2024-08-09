package instrumenter

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Instrumenter[REQUEST any, RESPONSE any] struct {
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
	propagator    propagation.TextMapPropagator
	carrierGetter func(REQUEST) propagation.TextMapCarrier
	base          Instrumenter[REQUEST, RESPONSE]
}

type PropagatingFromUpstreamInstrumenter[REQUEST any, RESPONSE any] struct {
	propagator    propagation.TextMapPropagator
	carrierGetter func(REQUEST) propagation.TextMapCarrier
	base          Instrumenter[REQUEST, RESPONSE]
}

func (i *Instrumenter[REQUEST, RESPONSE]) Start(parentContext context.Context, request REQUEST) context.Context {
	if len(i.operationListeners) > 0 {
		for _, listener := range i.operationListeners {
			parentContext = listener.OnBeforeStart(parentContext, time.Now())
		}
	}
	// extract span name
	spanName := i.spanNameExtractor.Extract(request)
	spanKind := i.spanKindExtractor.Extract(request)
	newCtx, span := i.tracer.Start(parentContext, spanName, trace.WithSpanKind(spanKind))
	var attributes []attribute.KeyValue
	// extract span attributes
	for _, extractor := range i.attributesExtractors {
		attributes = extractor.OnStart(attributes, newCtx, request)
	}
	// execute context customizer hook
	for _, customizer := range i.contextCustomizers {
		newCtx = customizer.OnStart(newCtx, request, attributes)
	}
	if len(i.operationListeners) > 0 {
		for _, listener := range i.operationListeners {
			newCtx = listener.OnBeforeEnd(newCtx, attributes, time.Now())
		}
	}
	span.SetAttributes(attributes...)
	return i.spanSuppressor.StoreInContext(newCtx, spanKind, span)
}

func (i *Instrumenter[REQUEST, RESPONSE]) End(ctx context.Context, request REQUEST, response RESPONSE, err error) {
	if len(i.operationListeners) > 0 {
		for _, listener := range i.operationListeners {
			listener.OnAfterStart(ctx, time.Now())
		}
	}
	span := trace.SpanFromContext(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	var attributes []attribute.KeyValue
	// extract span attributes
	for _, extractor := range i.attributesExtractors {
		attributes = extractor.OnEnd(attributes, ctx, request, response, err)
	}
	i.spanStatusExtractor.Extract(span, request, response, err)
	span.SetAttributes(attributes...)
	span.End()
	if len(i.operationListeners) > 0 {
		for _, listener := range i.operationListeners {
			listener.OnAfterEnd(ctx, attributes, time.Now())
		}
	}
}

func (p *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]) Start(parentContext context.Context, request REQUEST) context.Context {
	newCtx := p.base.Start(parentContext, request)
	p.propagator.Inject(newCtx, p.carrierGetter(request))
	return newCtx
}

func (p *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]) End(ctx context.Context, request REQUEST, response RESPONSE, err error) {
	p.base.End(ctx, request, response, err)
}

func (p *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]) Start(parentContext context.Context, request REQUEST) context.Context {
	extracted := p.propagator.Extract(parentContext, p.carrierGetter(request))
	return p.base.Start(extracted, request)
}

func (p *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]) End(ctx context.Context, request REQUEST, response RESPONSE, err error) {
	p.base.End(ctx, request, response, err)
}
