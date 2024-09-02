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
package dubbo

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"strconv"
)

type DubboCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 DubboCommonAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	DubboGetter GETTER1
	NetGetter   GETTER2
	Converter   DubboStatusCodeConverter
}

func (h *DubboCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   "dubbo.request.method",
		Value: attribute.StringValue(h.DubboGetter.GetRequestMethod(request)),
	})
	return attributes
}

func (h *DubboCommonAttrsExtractor[REQUEST, RESPONSE, GETTER, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	statusCode := h.DubboGetter.GetDubboResponseStatusCode(request, response, err)
	code, _ := strconv.Atoi(statusCode)
	protocolName := h.NetGetter.GetNetworkProtocolName(request, response)
	protocolVersion := h.NetGetter.GetNetworkProtocolVersion(request, response)
	attributes = append(attributes, attribute.KeyValue{
		Key:   "dubbo.response.status_code",
		Value: attribute.IntValue(code),
	}, attribute.KeyValue{
		Key:   semconv.NetworkProtocolNameKey,
		Value: attribute.StringValue(protocolName),
	}, attribute.KeyValue{
		Key:   semconv.NetworkProtocolVersionKey,
		Value: attribute.StringValue(protocolVersion),
	})
	return attributes
}

type DubboClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 DubboClientAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	Base             DubboCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	NetworkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
}

func (h *DubboClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.Base.OnStart(attributes, parentContext, request)
	return attributes
}

func (h *DubboClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.Base.OnEnd(attributes, context, request, response, err)
	attributes = h.NetworkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}

type DubboServerAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 DubboServerAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	Base             DubboCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	NetworkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
}

func (h *DubboServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.Base.OnStart(attributes, parentContext, request)
	return attributes
}

func (h *DubboServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.Base.OnEnd(attributes, context, request, response, err)
	attributes = h.NetworkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}
