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
