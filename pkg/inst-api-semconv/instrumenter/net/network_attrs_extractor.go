package net

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"strings"
)

type NetworkAttrsExtractor[REQUEST any, RESPONSE any, GETTER NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	getter GETTER
}

func (i *NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	return attributes
}

func (i *NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.NetworkTransportKey,
		Value: attribute.StringValue(i.getter.GetNetworkTransport(request, response)),
	}, attribute.KeyValue{
		Key:   semconv.NetworkTypeKey,
		Value: attribute.StringValue(strings.ToLower(i.getter.GetNetworkType(request, response))),
	}, attribute.KeyValue{
		Key:   semconv.NetworkProtocolNameKey,
		Value: attribute.StringValue(strings.ToLower(i.getter.GetNetworkProtocolName(request, response))),
	}, attribute.KeyValue{
		Key:   semconv.NetworkProtocolVersionKey,
		Value: attribute.StringValue(strings.ToLower(i.getter.GetNetworkProtocolVersion(request, response))),
	}, attribute.KeyValue{
		Key:   semconv.NetworkLocalAddressKey,
		Value: attribute.StringValue(i.getter.GetNetworkLocalInetAddress(request, response)),
	}, attribute.KeyValue{
		Key:   semconv.NetworkPeerAddressKey,
		Value: attribute.StringValue(i.getter.GetNetworkPeerInetAddress(request, response)),
	})
	localPort := i.getter.GetNetworkLocalPort(request, response)
	if localPort > 0 {
		attributes = append(attributes, attribute.KeyValue{
			Key:   semconv.NetworkLocalPortKey,
			Value: attribute.IntValue(localPort),
		})
	}
	peerPort := i.getter.GetNetworkPeerPort(request, response)
	if peerPort > 0 {
		attributes = append(attributes, attribute.KeyValue{
			Key:   semconv.NetworkPeerPortKey,
			Value: attribute.IntValue(peerPort),
		})
	}
	return attributes
}

type UrlAttrsExtractor[REQUEST any, RESPONSE any, GETTER UrlAttrsGetter[REQUEST]] struct {
	getter GETTER
	// TODO: add scheme provider for extension
}

func (u *UrlAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.URLSchemeKey,
		Value: attribute.StringValue(u.getter.GetUrlScheme(request)),
	}, attribute.KeyValue{
		Key:   semconv.URLPathKey,
		Value: attribute.StringValue(u.getter.GetUrlPath(request)),
	}, attribute.KeyValue{
		Key:   semconv.URLQueryKey,
		Value: attribute.StringValue(u.getter.GetUrlQuery(request)),
	})
	return attributes
}

func (u *UrlAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return attributes
}
