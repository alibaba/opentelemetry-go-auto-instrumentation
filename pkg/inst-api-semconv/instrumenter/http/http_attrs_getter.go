package http

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"

type HttpCommonAttrsGetter[REQUEST any, RESPONSE any] interface {
	GetRequestMethod(request REQUEST) string
	GetHttpRequestHeader(request REQUEST, name string) []string
	GetHttpResponseStatusCode(request REQUEST, response RESPONSE, err error) int
	GetHttpResponseHeader(request REQUEST, response RESPONSE, name string) []string
	GetErrorType(request REQUEST, response RESPONSE, err error) string
}

type HttpServerAttrsGetter[REQUEST any, RESPONSE any] interface {
	HttpCommonAttrsGetter[REQUEST, RESPONSE]
	net.UrlAttrsGetter[REQUEST]
	net.NetworkAttrsGetter[REQUEST, RESPONSE]
	GetUrlScheme(request REQUEST) string
	GetUrlPath(request REQUEST) string
	GetUrlQuery(request REQUEST) string
	GetHttpRoute(request REQUEST) string
}

type HttpClientAttrsGetter[REQUEST any, RESPONSE any] interface {
	HttpCommonAttrsGetter[REQUEST, RESPONSE]
	net.NetworkAttrsGetter[REQUEST, RESPONSE]
	GetUrlFull(request REQUEST) string
	GetServerAddress(request REQUEST) string
}
