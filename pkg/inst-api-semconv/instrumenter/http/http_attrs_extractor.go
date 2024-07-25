package http

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"go.opentelemetry.io/otel/attribute"
)

const http_request_method = attribute.Key("http.request.method")
const http_response_status_code = attribute.Key("http.response.status_code")
const http_route = attribute.Key("http.route")

const network_protocol_name = attribute.Key("network.protocol.name")
const network_protocol_version = attribute.Key("network.protocol.version")

const url_full = attribute.Key("url.full")

const user_agent_original = attribute.Key("user_agent.original")

// TODO: http.route

type HttpCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 HttpCommonAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	httpGetter GETTER1
	netGetter  GETTER2
	converter  HttpStatusCodeConverter
}

func (h *HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   http_request_method,
		Value: attribute.StringValue(h.httpGetter.GetRequestMethod(request)),
	})
	return attributes
}

func (h *HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	statusCode := h.httpGetter.GetHttpResponseStatusCode(request, response, err)
	protocolName := h.netGetter.GetNetworkProtocolName(request, response)
	protocolVersion := h.netGetter.GetNetworkProtocolVersion(request, response)
	attributes = append(attributes, attribute.KeyValue{
		Key:   http_response_status_code,
		Value: attribute.IntValue(statusCode),
	}, attribute.KeyValue{
		Key:   network_protocol_name,
		Value: attribute.StringValue(protocolName),
	}, attribute.KeyValue{
		Key:   network_protocol_version,
		Value: attribute.StringValue(protocolVersion),
	})
	return attributes
}

type HttpClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 HttpClientAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	base             HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	networkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
}

func (h *HttpClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.base.OnStart(attributes, parentContext, request)
	fullUrl := h.base.httpGetter.GetUrlFull(request)
	// TODO: add resend count
	attributes = append(attributes, attribute.KeyValue{
		Key:   url_full,
		Value: attribute.StringValue(fullUrl),
	})
	return attributes
}

func (h *HttpClientAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.base.OnEnd(attributes, context, request, response, err)
	attributes = h.networkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}

type HttpServerAttrsExtractor[REQUEST any, RESPONSE any, GETTER1 HttpServerAttrsGetter[REQUEST, RESPONSE], GETTER2 net.NetworkAttrsGetter[REQUEST, RESPONSE], GETTER3 net.UrlAttrsGetter[REQUEST]] struct {
	base             HttpCommonAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2]
	networkExtractor net.NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER2]
	urlExtractor     net.UrlAttrsExtractor[REQUEST, RESPONSE, GETTER3]
}

func (h *HttpServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2, GETTER3]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = h.base.OnStart(attributes, parentContext, request)
	attributes = h.urlExtractor.OnStart(attributes, parentContext, request)
	userAgent := h.base.httpGetter.GetHttpRequestHeader(request, "User-Agent")
	var firstUserAgent string
	if len(userAgent) > 0 {
		firstUserAgent = userAgent[0]
	} else {
		firstUserAgent = ""
	}
	attributes = append(attributes, attribute.KeyValue{Key: http_route,
		Value: attribute.StringValue(h.base.httpGetter.GetHttpRoute(request)),
	}, attribute.KeyValue{
		Key:   user_agent_original,
		Value: attribute.StringValue(firstUserAgent),
	})
	return attributes
}

func (h *HttpServerAttrsExtractor[REQUEST, RESPONSE, GETTER1, GETTER2, GETTER3]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = h.base.OnEnd(attributes, context, request, response, err)
	attributes = h.networkExtractor.OnEnd(attributes, context, request, response, err)
	return attributes
}
