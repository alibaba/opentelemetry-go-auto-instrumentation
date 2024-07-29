package instrumenter

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"testing"
)

func TestNoopSpanSupressor(t *testing.T) {
	n := NoopSpanSuppressor{}
	ctx := context.Background()
	n.StoreInContext(ctx, trace.SpanKindClient, noop.Span{})
	if n.ShouldSuppress(ctx, trace.SpanKindClient) != false {
		t.Errorf("should not suppress span")
	}
}

func TestSpanKeySuppressor(t *testing.T) {
	s := SpanKeySuppressor{
		spanKeys: []attribute.Key{
			utils.RPC_CLIENT_KEY,
		},
	}
	ctx := context.Background()
	newCtx := s.StoreInContext(ctx, trace.SpanKindClient, noop.Span{})
	if !s.ShouldSuppress(newCtx, trace.SpanKindClient) {
		t.Errorf("should suppress span")
	}
	if s.ShouldSuppress(context.Background(), trace.SpanKindClient) {
		t.Errorf("should not suppress span")
	}
}
