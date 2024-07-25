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

type alwaysInternalExtractor[REQUEST any] struct {
}

func (a *alwaysInternalExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindInternal
}

type alwaysClientExtractor[REQUEST any] struct {
}

func (a *alwaysClientExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindClient
}

type alwaysServerExtractor[REQUEST any] struct {
}

func (a *alwaysServerExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindServer
}

type alwaysProducerExtractor[REQUEST any] struct {
}

func (a *alwaysProducerExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindProducer
}

type alwaysConsumerExtractor[REQUEST any] struct {
}

func (a *alwaysConsumerExtractor[any]) Extract(request any) trace.SpanKind {
	return trace.SpanKindConsumer
}

func AlwaysInternalExtractor() SpanKindExtractor[any] {
	return &alwaysInternalExtractor[any]{}
}

func AlwaysClientExtractor() SpanKindExtractor[any] {
	return &alwaysClientExtractor[any]{}
}

func AlwaysServerExtractor() SpanKindExtractor[any] {
	return &alwaysServerExtractor[any]{}
}

func AlwaysProducerExtractor() SpanKindExtractor[any] {
	return &alwaysProducerExtractor[any]{}
}

func AlwaysConsumerExtractor() SpanKindExtractor[any] {
	return &alwaysConsumerExtractor[any]{}
}

type defaultSpanStatusExtractor[REQUEST any, RESPONSE any] struct {
}

func (d *defaultSpanStatusExtractor[REQUEST, RESPONSE]) Extract(span trace.Span, request REQUEST, response RESPONSE, err error) {
	if err != nil {
		span.SetStatus(codes.Error, "")
	}
}
