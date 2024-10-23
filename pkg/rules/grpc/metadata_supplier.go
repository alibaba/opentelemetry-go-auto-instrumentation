// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

type grpcMetadataSupplier struct {
	metadata *metadata.MD
}

// assert that grpcMetadataSupplier implements the TextMapCarrier interface.
var _ propagation.TextMapCarrier = &grpcMetadataSupplier{}

func (s *grpcMetadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *grpcMetadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *grpcMetadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}

func inject(ctx context.Context, propagators propagation.TextMapPropagator, methodName string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	propagators.Inject(ctx, &grpcMetadataSupplier{
		metadata: &md,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func extract(ctx context.Context, propagators propagation.TextMapPropagator) (context.Context, metadata.MD) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	return propagators.Extract(ctx, &grpcMetadataSupplier{
		metadata: &md,
	}), md
}
