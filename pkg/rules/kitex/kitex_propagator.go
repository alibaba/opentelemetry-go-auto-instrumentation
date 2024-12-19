package kitex

import (
	"context"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = &metadataProvider{}

type metadataProvider struct {
	metadata map[string]string
}

func (m *metadataProvider) Get(key string) string {
	if v, ok := m.metadata[key]; ok {
		return v
	}
	return ""
}

func (m *metadataProvider) Set(key, value string) {
	m.metadata[key] = value
}

func (m *metadataProvider) Keys() []string {
	out := make([]string, 0, len(m.metadata))
	for k := range m.metadata {
		out = append(out, k)
	}
	return out
}

func Inject(ctx context.Context, metadata map[string]string) {
	otel.GetTextMapPropagator().Inject(ctx, &metadataProvider{metadata: metadata})
}

func Extract(ctx context.Context, metadata map[string]string) context.Context {
	ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataProvider{metadata: CGIVariableToHTTPHeaderMetadata(metadata)})
	return ctx
}

func CGIVariableToHTTPHeaderMetadata(metadata map[string]string) map[string]string {
	res := make(map[string]string, len(metadata))
	for k, v := range metadata {
		res[metainfo.CGIVariableToHTTPHeader(k)] = v
	}
	return res
}
