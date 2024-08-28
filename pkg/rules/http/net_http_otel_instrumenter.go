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
//go:build ignore

package rule

import (
	"strconv"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/http"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type netHttpClientAttrsGetter struct {
}

func (n netHttpClientAttrsGetter) GetRequestMethod(request netHttpRequest) string {
	return request.method
}

func (n netHttpClientAttrsGetter) GetHttpRequestHeader(request netHttpRequest, name string) []string {
	return request.header.Values(name)
}

func (n netHttpClientAttrsGetter) GetHttpResponseStatusCode(request netHttpRequest, response netHttpResponse, err error) int {
	return response.statusCode
}

func (n netHttpClientAttrsGetter) GetHttpResponseHeader(request netHttpRequest, response netHttpResponse, name string) []string {
	return response.header.Values(name)
}

func (n netHttpClientAttrsGetter) GetErrorType(request netHttpRequest, response netHttpResponse, err error) string {
	return ""
}

func (n netHttpClientAttrsGetter) GetNetworkType(request netHttpRequest, response netHttpResponse) string {
	return "ipv4"
}

func (n netHttpClientAttrsGetter) GetNetworkTransport(request netHttpRequest, response netHttpResponse) string {
	return "tcp"
}

func (n netHttpClientAttrsGetter) GetNetworkProtocolName(request netHttpRequest, response netHttpResponse) string {
	if request.isTls == false {
		return "http"
	} else {
		return "https"
	}
}

func (n netHttpClientAttrsGetter) GetNetworkProtocolVersion(request netHttpRequest, response netHttpResponse) string {
	return request.version
}

func (n netHttpClientAttrsGetter) GetNetworkLocalInetAddress(request netHttpRequest, response netHttpResponse) string {
	return ""
}

func (n netHttpClientAttrsGetter) GetNetworkLocalPort(request netHttpRequest, response netHttpResponse) int {
	return 0
}

func (n netHttpClientAttrsGetter) GetNetworkPeerInetAddress(request netHttpRequest, response netHttpResponse) string {
	return request.host
}

func (n netHttpClientAttrsGetter) GetNetworkPeerPort(request netHttpRequest, response netHttpResponse) int {
	port, err := strconv.Atoi(request.url.Port())
	if err != nil {
		return 0
	}
	return port
}

func (n netHttpClientAttrsGetter) GetUrlFull(request netHttpRequest) string {
	return request.url.String()
}

func (n netHttpClientAttrsGetter) GetServerAddress(request netHttpRequest) string {
	return request.host
}

type netHttpServerAttrsGetter struct {
}

func (n netHttpServerAttrsGetter) GetRequestMethod(request netHttpRequest) string {
	return request.method
}

func (n netHttpServerAttrsGetter) GetHttpRequestHeader(request netHttpRequest, name string) []string {
	return request.header.Values(name)
}

func (n netHttpServerAttrsGetter) GetHttpResponseStatusCode(request netHttpRequest, response netHttpResponse, err error) int {
	return response.statusCode
}

func (n netHttpServerAttrsGetter) GetHttpResponseHeader(request netHttpRequest, response netHttpResponse, name string) []string {
	return response.header.Values(name)
}

func (n netHttpServerAttrsGetter) GetErrorType(request netHttpRequest, response netHttpResponse, err error) string {
	return ""
}

func (n netHttpServerAttrsGetter) GetUrlScheme(request netHttpRequest) string {
	return request.url.Scheme
}

func (n netHttpServerAttrsGetter) GetUrlPath(request netHttpRequest) string {
	return request.url.Path
}

func (n netHttpServerAttrsGetter) GetUrlQuery(request netHttpRequest) string {
	return request.url.RawQuery
}

func (n netHttpServerAttrsGetter) GetNetworkType(request netHttpRequest, response netHttpResponse) string {
	return "ipv4"
}

func (n netHttpServerAttrsGetter) GetNetworkTransport(request netHttpRequest, response netHttpResponse) string {
	return "tcp"
}

func (n netHttpServerAttrsGetter) GetNetworkProtocolName(request netHttpRequest, response netHttpResponse) string {
	if request.isTls == false {
		return "http"
	} else {
		return "https"
	}
}

func (n netHttpServerAttrsGetter) GetNetworkProtocolVersion(request netHttpRequest, response netHttpResponse) string {
	return request.version
}

func (n netHttpServerAttrsGetter) GetNetworkLocalInetAddress(request netHttpRequest, response netHttpResponse) string {
	return ""
}

func (n netHttpServerAttrsGetter) GetNetworkLocalPort(request netHttpRequest, response netHttpResponse) int {
	return 0
}

func (n netHttpServerAttrsGetter) GetNetworkPeerInetAddress(request netHttpRequest, response netHttpResponse) string {
	return request.host
}

func (n netHttpServerAttrsGetter) GetNetworkPeerPort(request netHttpRequest, response netHttpResponse) int {
	port, err := strconv.Atoi(request.url.Port())
	if err != nil {
		return 0
	}
	return port
}

func (n netHttpServerAttrsGetter) GetHttpRoute(request netHttpRequest) string {
	return request.url.Path
}

func BuildNetHttpClientOtelInstrumenter() *instrumenter.PropagatingToDownstreamInstrumenter[netHttpRequest, netHttpResponse] {
	builder := instrumenter.Builder[netHttpRequest, netHttpResponse]{}
	clientGetter := netHttpClientAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[netHttpRequest, netHttpResponse, netHttpClientAttrsGetter, netHttpClientAttrsGetter]{HttpGetter: clientGetter, NetGetter: clientGetter, Converter: &http.ClientHttpStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[netHttpRequest, netHttpResponse, netHttpClientAttrsGetter]{Getter: clientGetter}
	return builder.Init().SetSpanNameExtractor(&http.HttpClientSpanNameExtractor[netHttpRequest, netHttpResponse]{Getter: clientGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[netHttpRequest]{}).
		AddAttributesExtractor(&http.HttpClientAttrsExtractor[netHttpRequest, netHttpResponse, netHttpClientAttrsGetter, netHttpClientAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor}).BuildPropagatingToDownstreamInstrumenter(func(n netHttpRequest) propagation.TextMapCarrier {
		if n.header == nil {
			return nil
		}
		return propagation.HeaderCarrier(n.header)
	}, otel.GetTextMapPropagator())
}

func BuildNetHttpServerOtelInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[netHttpRequest, netHttpResponse] {
	builder := instrumenter.Builder[netHttpRequest, netHttpResponse]{}
	serverGetter := netHttpServerAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[netHttpRequest, netHttpResponse, netHttpServerAttrsGetter, netHttpServerAttrsGetter]{HttpGetter: serverGetter, NetGetter: serverGetter, Converter: &http.ServerHttpStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[netHttpRequest, netHttpResponse, netHttpServerAttrsGetter]{Getter: serverGetter}
	urlExtractor := net.UrlAttrsExtractor[netHttpRequest, netHttpResponse, netHttpServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanNameExtractor(&http.HttpServerSpanNameExtractor[netHttpRequest, netHttpResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[netHttpRequest]{}).
		AddAttributesExtractor(&http.HttpServerAttrsExtractor[netHttpRequest, netHttpResponse, netHttpServerAttrsGetter, netHttpServerAttrsGetter, netHttpServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor, UrlExtractor: urlExtractor}).BuildPropagatingFromUpstreamInstrumenter(func(n netHttpRequest) propagation.TextMapCarrier {
		if n.header == nil {
			return nil
		}
		return propagation.HeaderCarrier(n.header)
	}, otel.GetTextMapPropagator())
}
