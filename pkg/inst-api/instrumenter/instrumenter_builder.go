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
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/core/meter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TODO: add route updater here, now we do not support such controller layer to update route.
type OperationMetrics interface {
	Create(meter metric.Meter) *OperationListenerWrapper
	Match(meta meter.MeterMeta) bool
}

type InstrumentEnabler interface {
	IsEnabled() bool
}

type defaultInstrumentEnabler struct {
}

func NewDefaultInstrumentEnabler() InstrumentEnabler {
	return &defaultInstrumentEnabler{}
}

func (a *defaultInstrumentEnabler) IsEnabled() bool {
	return true
}

type Builder[REQUEST any, RESPONSE any] struct {
	Enabler              InstrumentEnabler
	SpanNameExtractor    SpanNameExtractor[REQUEST]
	SpanKindExtractor    SpanKindExtractor[REQUEST]
	SpanStatusExtractor  SpanStatusExtractor[REQUEST, RESPONSE]
	AttributesExtractors []AttributesExtractor[REQUEST, RESPONSE]
	OperationListeners   []*OperationListenerWrapper
	OperationMetrics     []OperationMetrics
	ContextCustomizers   []ContextCustomizer[REQUEST]
	SpanSuppressor       SpanSuppressor
	Tracer               trace.Tracer
	InstVersion          string
}

func (b *Builder[REQUEST, RESPONSE]) Init() *Builder[REQUEST, RESPONSE] {
	b.Enabler = &defaultInstrumentEnabler{}
	b.AttributesExtractors = make([]AttributesExtractor[REQUEST, RESPONSE], 0)
	b.OperationListeners = b.buildOperationListeners()
	b.ContextCustomizers = make([]ContextCustomizer[REQUEST], 0)
	b.SpanSuppressor = b.buildSpanSuppressor()
	b.SpanStatusExtractor = &defaultSpanStatusExtractor[REQUEST, RESPONSE]{}
	b.Tracer = otel.GetTracerProvider().Tracer("")
	return b
}

func (b *Builder[REQUEST, RESPONSE]) SetInstrumentEnabler(enabler InstrumentEnabler) *Builder[REQUEST, RESPONSE] {
	b.Enabler = enabler
	return b
}

func (b *Builder[REQUEST, RESPONSE]) SetInstVersion(instVersion string) *Builder[REQUEST, RESPONSE] {
	b.InstVersion = instVersion
	return b
}

func (b *Builder[REQUEST, RESPONSE]) SetSpanNameExtractor(spanNameExtractor SpanNameExtractor[REQUEST]) *Builder[REQUEST, RESPONSE] {
	b.SpanNameExtractor = spanNameExtractor
	return b
}

func (b *Builder[REQUEST, RESPONSE]) SetSpanStatusExtractor(spanStatusExtractor SpanStatusExtractor[REQUEST, RESPONSE]) *Builder[REQUEST, RESPONSE] {
	b.SpanStatusExtractor = spanStatusExtractor
	return b
}

func (b *Builder[REQUEST, RESPONSE]) SetSpanKindExtractor(spanKindExtractor SpanKindExtractor[REQUEST]) *Builder[REQUEST, RESPONSE] {
	b.SpanKindExtractor = spanKindExtractor
	return b
}

func (b *Builder[REQUEST, RESPONSE]) AddAttributesExtractor(attributesExtractor ...AttributesExtractor[REQUEST, RESPONSE]) *Builder[REQUEST, RESPONSE] {
	b.AttributesExtractors = append(b.AttributesExtractors, attributesExtractor...)
	return b
}

func (b *Builder[REQUEST, RESPONSE]) AddOperationListeners(operationListener ...*OperationListenerWrapper) *Builder[REQUEST, RESPONSE] {
	b.OperationListeners = append(b.OperationListeners, operationListener...)
	return b
}

func (b *Builder[REQUEST, RESPONSE]) AddContextCustomizers(contextCustomizers ...ContextCustomizer[REQUEST]) *Builder[REQUEST, RESPONSE] {
	b.ContextCustomizers = append(b.ContextCustomizers, contextCustomizers...)
	return b
}

func (b *Builder[REQUEST, RESPONSE]) BuildInstrumenter() *InternalInstrumenter[REQUEST, RESPONSE] {
	return &InternalInstrumenter[REQUEST, RESPONSE]{
		enabler:              b.Enabler,
		spanNameExtractor:    b.SpanNameExtractor,
		spanKindExtractor:    b.SpanKindExtractor,
		spanStatusExtractor:  b.SpanStatusExtractor,
		attributesExtractors: b.AttributesExtractors,
		operationListeners:   b.OperationListeners,
		operationMetrics:     b.OperationMetrics,
		contextCustomizers:   b.ContextCustomizers,
		spanSuppressor:       b.SpanSuppressor,
		tracer:               b.Tracer,
		instVersion:          b.InstVersion,
	}
}

func (b *Builder[REQUEST, RESPONSE]) BuildPropagatingToDownstreamInstrumenter(carrierGetter func(REQUEST) propagation.TextMapCarrier, prop propagation.TextMapPropagator) *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE] {
	return &PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]{
		base: InternalInstrumenter[REQUEST, RESPONSE]{
			enabler:              b.Enabler,
			spanNameExtractor:    b.SpanNameExtractor,
			spanKindExtractor:    b.SpanKindExtractor,
			spanStatusExtractor:  b.SpanStatusExtractor,
			attributesExtractors: b.AttributesExtractors,
			operationListeners:   b.OperationListeners,
			operationMetrics:     b.OperationMetrics,
			contextCustomizers:   b.ContextCustomizers,
			spanSuppressor:       b.SpanSuppressor,
			tracer:               b.Tracer,
			instVersion:          b.InstVersion,
		},
		carrierGetter: carrierGetter,
		prop:          prop,
	}
}

func (b *Builder[REQUEST, RESPONSE]) BuildPropagatingToDownstreamInstrumenterWithProp(carrierGetter func(REQUEST) propagation.TextMapCarrier, prop propagation.TextMapPropagator) *PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE] {
	return &PropagatingToDownstreamInstrumenter[REQUEST, RESPONSE]{
		base: InternalInstrumenter[REQUEST, RESPONSE]{
			enabler:              b.Enabler,
			spanNameExtractor:    b.SpanNameExtractor,
			spanKindExtractor:    b.SpanKindExtractor,
			spanStatusExtractor:  b.SpanStatusExtractor,
			attributesExtractors: b.AttributesExtractors,
			operationListeners:   b.OperationListeners,
			operationMetrics:     b.OperationMetrics,
			contextCustomizers:   b.ContextCustomizers,
			spanSuppressor:       b.SpanSuppressor,
			tracer:               b.Tracer,
			instVersion:          b.InstVersion,
		},
		carrierGetter: carrierGetter,
		prop:          prop,
	}
}

func (b *Builder[REQUEST, RESPONSE]) BuildPropagatingFromUpstreamInstrumenterWithProp(carrierGetter func(REQUEST) propagation.TextMapCarrier, prop propagation.TextMapPropagator) *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE] {
	return &PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]{
		base: InternalInstrumenter[REQUEST, RESPONSE]{
			enabler:              b.Enabler,
			spanNameExtractor:    b.SpanNameExtractor,
			spanKindExtractor:    b.SpanKindExtractor,
			spanStatusExtractor:  b.SpanStatusExtractor,
			attributesExtractors: b.AttributesExtractors,
			operationListeners:   b.OperationListeners,
			operationMetrics:     b.OperationMetrics,
			contextCustomizers:   b.ContextCustomizers,
			spanSuppressor:       b.SpanSuppressor,
			tracer:               b.Tracer,
			instVersion:          b.InstVersion,
		},
		carrierGetter: carrierGetter,
		prop:          prop,
	}
}

func (b *Builder[REQUEST, RESPONSE]) BuildPropagatingFromUpstreamInstrumenter(carrierGetter func(REQUEST) propagation.TextMapCarrier, prop propagation.TextMapPropagator) *PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE] {
	return &PropagatingFromUpstreamInstrumenter[REQUEST, RESPONSE]{
		base: InternalInstrumenter[REQUEST, RESPONSE]{
			enabler:              b.Enabler,
			spanNameExtractor:    b.SpanNameExtractor,
			spanKindExtractor:    b.SpanKindExtractor,
			spanStatusExtractor:  b.SpanStatusExtractor,
			attributesExtractors: b.AttributesExtractors,
			operationListeners:   b.OperationListeners,
			operationMetrics:     b.OperationMetrics,
			contextCustomizers:   b.ContextCustomizers,
			spanSuppressor:       b.SpanSuppressor,
			tracer:               b.Tracer,
			instVersion:          b.InstVersion,
		},
		carrierGetter: carrierGetter,
		prop:          prop,
	}
}

func (b *Builder[REQUEST, RESPONSE]) buildOperationListeners() []*OperationListenerWrapper {
	if len(b.OperationMetrics) == 0 {
		return make([]*OperationListenerWrapper, 0)
	}
	meterProvider := meter.GetMeterProvider()
	if meterProvider == nil {
		return make([]*OperationListenerWrapper, 0)
	}

	listeners := make([]*OperationListenerWrapper, 0, len(b.OperationMetrics)+len(b.OperationListeners))

	meters := meterProvider.GetMeters()
	for _, m := range meters {
		for _, factory := range b.OperationMetrics {
			if factory.Match(m.Metas()) {
				listeners = append(listeners, factory.Create(m.Meter()))
			}
		}
	}
	return listeners
}

// TODO: create suppressor by otel.instrumentation.experimental.span-suppression-strategy
func (b *Builder[REQUEST, RESPONSE]) buildSpanSuppressor() SpanSuppressor {
	kvs := make(map[attribute.Key]bool)
	for _, extractor := range b.AttributesExtractors {
		provider, ok := extractor.(SpanKeyProvider)
		if ok {
			kvs[provider.GetSpanKey()] = true
		}
	}
	kSlice := make([]attribute.Key, 0, len(kvs))
	for k := range kvs {
		kSlice = append(kSlice, k)
	}
	return NewSpanKeySuppressor(kSlice)
}
