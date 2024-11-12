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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
	"sync"
)

var noneSuppressor *NoopSpanSuppressor
var bySpanKeySuppressor *SpanKeySuppressor
var spanKindSuppressor *SpanKindSuppressor
var once sync.Once

type SpanSuppressorStrategy interface {
	create(spanKeys []attribute.Key) SpanSuppressor
}

type SemConvStrategy struct{}

func (t *SemConvStrategy) create(spanKeys []attribute.Key) SpanSuppressor {
	once.Do(func() {
		bySpanKeySuppressor = NewSpanKeySuppressor(spanKeys)
	})
	return bySpanKeySuppressor
}

type NoneStrategy struct{}

func (n *NoneStrategy) create(spanKeys []attribute.Key) SpanSuppressor {
	once.Do(func() {
		noneSuppressor = NewNoopSpanSuppressor()
	})
	return noneSuppressor
}

type SpanKindStrategy struct{}

func (s *SpanKindStrategy) create(spanKeys []attribute.Key) SpanSuppressor {
	once.Do(func() {
		spanKindSuppressor = NewSpanKindSuppressor()
	})
	return spanKindSuppressor
}

type SpanKindSuppressor struct {
	delegates map[trace.SpanKind]SpanSuppressor
}

func getSpanSuppressionStrategyFromEnv() SpanSuppressorStrategy {
	suppressionStrategy := os.Getenv("OTEL_INSTRUMENTATION_EXPERIMENTAL_SPAN_SUPPRESSION_STRATEGY")
	switch suppressionStrategy {
	case "none":
		return &NoneStrategy{}
	case "span-kind":
		return &SpanKindStrategy{}
	default:
		return &SemConvStrategy{}
	}
}
