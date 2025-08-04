package sentinel

import (
	"context"
	"time"

	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/utils"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type sentinelInstrumenter struct {
	tracer trace.Tracer
}

func NewSentinelInstrumenter() *sentinelInstrumenter {
	tracer := otel.GetTracerProvider().
		Tracer(utils.SENTINEL_SCOPE_NAME,
			trace.WithInstrumentationVersion(version.Tag),
		)
	return &sentinelInstrumenter{
		tracer: tracer,
	}
}

func (s *sentinelInstrumenter) StartAndEnd(ctx context.Context, spanName string, StartTime time.Time, EndTime time.Time, attrs []attribute.KeyValue, opts ...trace.SpanStartOption) {
	// start and end span
	_, span := s.tracer.Start(context.Background(),
		spanName,
		trace.WithAttributes(attrs...),
		trace.WithTimestamp(StartTime),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	span.End(trace.WithTimestamp(EndTime))
}
