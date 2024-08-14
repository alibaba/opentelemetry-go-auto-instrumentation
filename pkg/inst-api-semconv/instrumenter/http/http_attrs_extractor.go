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
package http

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// TODO: http.route

type HttpCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 HttpCommonAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	HttpGetter GETTER1
	NetGetter  GETTER2
	Converter  HttpStatusCodeConverter
}

func (h *HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.HTTPRequestMethodKey,
		Value: attribute.StringValue(h.HttpGetter.GetRequestMethod(request)),
	})
	return attributes
}

func (h *HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	statusCode := h.HttpGetter.GetHttpResponseStatusCode(request, response, err)
	protocolName := h.NetGetter.GetNetworkProtocolName(request, response)
	protocolVersion := h.NetGetter.GetNetworkProtocolVersion(request, response)
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.HTTPResponseStatusCodeKey,
		Value: attribute.IntValue(statusCode),
	}, attribute.KeyValue{
		Key:   semconv.NetworkProtocolNameKey,
		Value: attribute.StringValue(protocolName),
	}, attribute.KeyValue{
		Key:   semconv.NetworkProtocolVersionKey,
		Value: attribute.StringValue(protocolVersion),
	})
	return attributes
}

type HttpClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 HttpClientAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	Base             HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	NetworkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
}

func (h *HttpClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.Base.OnStart(attributes, parentContext, request)
	fullUrl := h.Base.HttpGetter.GetUrlFull(request)
	// TODO: add resend count
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.URLFullKey,
		Value: attribute.StringValue(fullUrl),
	})
	return attributes
}

func (h *HttpClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.Base.OnEnd(attributes, context, request, response, err)
	attributes = h.NetworkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}

type HttpServerAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 HttpServerAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE], GETTER3 net.UrlAttrsGetter[REQUEST]] struct {
	Base             HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	NetworkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
	UrlExtractor     net.UrlAttrsExtractor[REQUEST, RESPONSE, GETTER3]
}

func (h *HttpServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2, GETTER3]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.Base.OnStart(attributes, parentContext, request)
	attributes = h.UrlExtractor.OnStart(attributes, parentContext, request)
	userAgent := h.Base.HttpGetter.GetHttpRequestHeader(request, "User-Agent")
	var firstUserAgent string
	if len(userAgent) > 0 {
		firstUserAgent = userAgent[0]
	} else {
		firstUserAgent = ""
	}
	attributes = append(attributes, attribute.KeyValue{Key: semconv.HTTPRouteKey,
		Value: attribute.StringValue(h.Base.HttpGetter.GetHttpRoute(request)),
	}, attribute.KeyValue{
		Key:   semconv.UserAgentOriginalKey,
		Value: attribute.StringValue(firstUserAgent),
	})
	return attributes
}

func (h *HttpServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2, GETTER3]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.Base.OnEnd(attributes, context, request, response, err)
	attributes = h.NetworkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}
