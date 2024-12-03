// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package instrumenter

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
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

type ContextCustomizer[REQUEST interface{}] interface {
	OnStart(context context.Context, request REQUEST, startAttributes []attribute.KeyValue) context.Context
}
