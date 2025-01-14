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
