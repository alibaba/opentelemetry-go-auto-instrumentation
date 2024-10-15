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

type ContextCustomizer[REQUEST interface{}] interface {
	OnStart(context context.Context, request REQUEST, startAttributes []attribute.KeyValue) context.Context
}
