package instrumenter

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type AttributesExtractor[REQUEST any, RESPONSE any] interface {
	OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue
	OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue
}

type SpanKindExtractor[REQUEST any] interface {
	Extract(request REQUEST) trace.SpanKind
}

type SpanNameExtractor[REQUEST any] interface {
	Extract(request REQUEST) string
}

type SpanStatusExtractor[REQUEST any, RESPONSE any] interface {
	Extract(span trace.Span, request REQUEST, response RESPONSE, err error)
}

type SpanKeyProvider interface {
	GetSpanKey() attribute.Key
}

type AlwaysInternalExtractor[REQUEST any] struct {
}

func (a *AlwaysInternalExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindInternal
}

type AlwaysClientExtractor[REQUEST any] struct {
}

func (a *AlwaysClientExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindClient
}

type AlwaysServerExtractor[REQUEST any] struct {
}

func (a *AlwaysServerExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindServer
}

type AlwaysProducerExtractor[REQUEST any] struct {
}

func (a *AlwaysProducerExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindProducer
}

type AlwaysConsumerExtractor[REQUEST any] struct {
}

func (a *AlwaysConsumerExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindConsumer
}

type defaultSpanStatusExtractor[REQUEST any, RESPONSE any] struct {
}

func (d *defaultSpanStatusExtractor[REQUEST, RESPONSE]) Extract(span trace.Span, request REQUEST, response RESPONSE, err error) {
	if err != nil {
		span.SetStatus(codes.Error, "")
	}
}
