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
//go:build ignore

package rule

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// TODO: http.route

type GrpcCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 GrpcCommonAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	GrpcGetter GETTER1
	Converter  GrpcStatusCodeConverter
}

func (h *GrpcCommonAttrsExtractor[REQUEST, RESPONSE, GETTER, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   "grpc.request.method",
		Value: attribute.StringValue(h.GrpcGetter.GetRequestMethod(request)),
	})
	return attributes
}

func (h *GrpcCommonAttrsExtractor[REQUEST, RESPONSE, GETTER, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	statusCode := h.GrpcGetter.GetGrpcResponseStatusCode(request, response, err)
	attributes = append(attributes, attribute.KeyValue{
		Key:   "grpc.response.status_code",
		Value: attribute.IntValue(statusCode),
	})
	return attributes
}

type GrpcClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 GrpcClientAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	Base             GrpcCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	NetworkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
}

func (h *GrpcClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.Base.OnStart(attributes, parentContext, request)
	fullUrl := h.Base.GrpcGetter.GetUrlFull(request)
	// TODO: add resend count
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.URLFullKey,
		Value: attribute.StringValue(fullUrl),
	})
	return attributes
}

func (h *GrpcClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.Base.OnEnd(attributes, context, request, response, err)
	attributes = h.NetworkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}

type GrpcServerAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 GrpcServerAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	Base             GrpcCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	NetworkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
}

func (h *GrpcServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.Base.OnStart(attributes, parentContext, request)
	return attributes
}

func (h *GrpcServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.Base.OnEnd(attributes, context, request, response, err)
	attributes = h.NetworkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}
