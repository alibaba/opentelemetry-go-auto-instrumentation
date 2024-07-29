package instrumenter

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

type OperationListener interface {
	OnBeforeStart(parentContext context.Context, startTimestamp time.Time) context.Context
	OnBeforeEnd(context context.Context, startAttributes []attribute.KeyValue, startTimestamp time.Time) context.Context
	OnAfterStart(context context.Context, endTimestamp time.Time)
	OnAfterEnd(context context.Context, endAttributes []attribute.KeyValue, endTimestamp time.Time)
}

type AttrsShadower interface {
	Shadow(attrs []attribute.KeyValue) (int, []attribute.KeyValue)
}

type NoopAttrsShadower struct{}

func (n NoopAttrsShadower) Shadow(attrs []attribute.KeyValue) (int, []attribute.KeyValue) {
	return len(attrs), attrs
}

type OperationListenerWrapper struct {
	listener       OperationListener
	attrCustomizer AttrsShadower
}

func (w *OperationListenerWrapper) OnBeforeStart(parentContext context.Context, startTimestamp time.Time) context.Context {
	return w.listener.OnBeforeStart(parentContext, startTimestamp)
}

func (w *OperationListenerWrapper) OnBeforeEnd(context context.Context, startAttributes []attribute.KeyValue, startTimestamp time.Time) context.Context {
	if w.attrCustomizer != nil {
		validNum, startAttributes := w.attrCustomizer.Shadow(startAttributes)
		return w.listener.OnBeforeEnd(context, startAttributes[:validNum], startTimestamp)
	} else {
		return w.listener.OnBeforeEnd(context, startAttributes, startTimestamp)
	}
}

func (w *OperationListenerWrapper) OnAfterStart(context context.Context, endTimestamp time.Time) {
	w.listener.OnAfterStart(context, endTimestamp)
}

func (w *OperationListenerWrapper) OnAfterEnd(context context.Context, endAttributes []attribute.KeyValue, endTimestamp time.Time) {
	if w.attrCustomizer != nil {
		validNum, endAttributes := w.attrCustomizer.Shadow(endAttributes)
		w.listener.OnAfterEnd(context, endAttributes[:validNum], endTimestamp)
	} else {
		w.listener.OnAfterEnd(context, endAttributes, endTimestamp)
	}
}

type ContextCustomizer[REQUEST interface{}] interface {
	OnStart(context context.Context, request REQUEST, startAttributes []attribute.KeyValue) context.Context
}
