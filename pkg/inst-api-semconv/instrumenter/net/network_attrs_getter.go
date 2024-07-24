package net

type NetworkAttrsGetter[REQUEST any, RESPONSE any] interface {
	GetNetworkType(request REQUEST, response RESPONSE) string
	GetNetworkTransport(request REQUEST, response RESPONSE) string
	GetNetworkProtocolName(request REQUEST, response RESPONSE) string
	GetNetworkProtocolVersion(request REQUEST, response RESPONSE) string
	GetNetworkLocalInetAddress(request REQUEST, response RESPONSE) string
	GetNetworkLocalPort(request REQUEST, response RESPONSE) int
	GetNetworkPeerInetAddress(request REQUEST, response RESPONSE) string
	GetNetworkPeerPort(request REQUEST, response RESPONSE) int
}

type UrlAttrsGetter[REQUEST any] interface {
	GetUrlScheme(request REQUEST) string
	GetUrlPath(request REQUEST) string
	GetUrlQuery(request REQUEST) string
}
