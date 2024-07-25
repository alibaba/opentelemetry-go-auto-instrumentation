package net

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"strings"
)

const network_transport = attribute.Key("network.transport")
const network_protocol_name = attribute.Key("network.protocol.name")
const network_local_address = attribute.Key("network.local.address")
const network_local_port = attribute.Key("network.local.port")
const network_peer_address = attribute.Key("network.peer.address")
const network_peer_port = attribute.Key("network.peer.port")
const network_protocol_version = attribute.Key("network.protocol.version")
const network_type = attribute.Key("network.type")

const url_scheme = attribute.Key("url.scheme")
const url_query = attribute.Key("url.query")
const url_path = attribute.Key("url.path")

type NetworkAttrsExtractor[REQUEST any, RESPONSE any, GETTER NetworkAttrsGetter[REQUEST, RESPONSE]] struct {
	getter GETTER
}

func (i *NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	return attributes
}

func (i *NetworkAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   network_transport,
		Value: attribute.StringValue(i.getter.GetNetworkTransport(request, response)),
	}, attribute.KeyValue{
		Key:   network_type,
		Value: attribute.StringValue(strings.ToLower(i.getter.GetNetworkType(request, response))),
	}, attribute.KeyValue{
		Key:   network_protocol_name,
		Value: attribute.StringValue(strings.ToLower(i.getter.GetNetworkProtocolName(request, response))),
	}, attribute.KeyValue{
		Key:   network_protocol_version,
		Value: attribute.StringValue(strings.ToLower(i.getter.GetNetworkProtocolVersion(request, response))),
	}, attribute.KeyValue{
		Key:   network_local_address,
		Value: attribute.StringValue(i.getter.GetNetworkLocalInetAddress(request, response)),
	}, attribute.KeyValue{
		Key:   network_peer_address,
		Value: attribute.StringValue(i.getter.GetNetworkPeerInetAddress(request, response)),
	})
	localPort := i.getter.GetNetworkLocalPort(request, response)
	if localPort > 0 {
		attributes = append(attributes, attribute.KeyValue{
			Key:   network_local_port,
			Value: attribute.IntValue(localPort),
		})
	}
	peerPort := i.getter.GetNetworkPeerPort(request, response)
	if peerPort > 0 {
		attributes = append(attributes, attribute.KeyValue{
			Key:   network_peer_port,
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
		Key:   url_scheme,
		Value: attribute.StringValue(u.getter.GetUrlScheme(request)),
	}, attribute.KeyValue{
		Key:   url_path,
		Value: attribute.StringValue(u.getter.GetUrlPath(request)),
	}, attribute.KeyValue{
		Key:   url_query,
		Value: attribute.StringValue(u.getter.GetUrlQuery(request)),
	})
	return attributes
}

func (u *UrlAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return attributes
}
