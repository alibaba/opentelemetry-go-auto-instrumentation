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
