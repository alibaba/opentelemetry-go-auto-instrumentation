// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
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

type muxHttpServerAttrsGetter struct {
}

func (n muxHttpServerAttrsGetter) GetRequestMethod(request muxHttpRequest) string {
	return request.method
}

func (n muxHttpServerAttrsGetter) GetHttpRequestHeader(request muxHttpRequest, name string) []string {
	return request.header.Values(name)
}

func (n muxHttpServerAttrsGetter) GetHttpResponseStatusCode(request muxHttpRequest, response muxHttpResponse, err error) int {
	return response.statusCode
}

func (n muxHttpServerAttrsGetter) GetHttpResponseHeader(request muxHttpRequest, response muxHttpResponse, name string) []string {
	return response.header.Values(name)
}

func (n muxHttpServerAttrsGetter) GetErrorType(request muxHttpRequest, response muxHttpResponse, err error) string {
	return ""
}

func (n muxHttpServerAttrsGetter) GetUrlScheme(request muxHttpRequest) string {
	return request.url.Scheme
}

func (n muxHttpServerAttrsGetter) GetUrlPath(request muxHttpRequest) string {
	return request.url.Path
}

func (n muxHttpServerAttrsGetter) GetUrlQuery(request muxHttpRequest) string {
	return request.url.RawQuery
}

func (n muxHttpServerAttrsGetter) GetNetworkType(request muxHttpRequest, response muxHttpResponse) string {
	return "ipv4"
}

func (n muxHttpServerAttrsGetter) GetNetworkTransport(request muxHttpRequest, response muxHttpResponse) string {
	return "tcp"
}

func (n muxHttpServerAttrsGetter) GetNetworkProtocolName(request muxHttpRequest, response muxHttpResponse) string {
	if request.isTls == false {
		return "http"
	} else {
		return "https"
	}
}

func (n muxHttpServerAttrsGetter) GetNetworkProtocolVersion(request muxHttpRequest, response muxHttpResponse) string {
	return request.version
}

func (n muxHttpServerAttrsGetter) GetNetworkLocalInetAddress(request muxHttpRequest, response muxHttpResponse) string {
	return ""
}

func (n muxHttpServerAttrsGetter) GetNetworkLocalPort(request muxHttpRequest, response muxHttpResponse) int {
	return 0
}

func (n muxHttpServerAttrsGetter) GetNetworkPeerInetAddress(request muxHttpRequest, response muxHttpResponse) string {
	return request.host
}

func (n muxHttpServerAttrsGetter) GetNetworkPeerPort(request muxHttpRequest, response muxHttpResponse) int {
	port, err := strconv.Atoi(request.url.Port())
	if err != nil {
		return 0
	}
	return port
}

func (n muxHttpServerAttrsGetter) GetHttpRoute(request muxHttpRequest) string {
	return request.url.Path
}

func BuildMuxHttpServerOtelInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[muxHttpRequest, muxHttpResponse] {
	builder := instrumenter.Builder[muxHttpRequest, muxHttpResponse]{}
	serverGetter := muxHttpServerAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[muxHttpRequest, muxHttpResponse, muxHttpServerAttrsGetter, muxHttpServerAttrsGetter]{HttpGetter: serverGetter, NetGetter: serverGetter, Converter: &http.ServerHttpStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[muxHttpRequest, muxHttpResponse, muxHttpServerAttrsGetter]{Getter: serverGetter}
	urlExtractor := net.UrlAttrsExtractor[muxHttpRequest, muxHttpResponse, muxHttpServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanNameExtractor(&http.HttpServerSpanNameExtractor[muxHttpRequest, muxHttpResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[muxHttpRequest]{}).
		AddAttributesExtractor(&http.HttpServerAttrsExtractor[muxHttpRequest, muxHttpResponse, muxHttpServerAttrsGetter, muxHttpServerAttrsGetter, muxHttpServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor, UrlExtractor: urlExtractor}).BuildPropagatingFromUpstreamInstrumenter(func(n muxHttpRequest) propagation.TextMapCarrier {
		if n.header == nil {
			return nil
		}
		return propagation.HeaderCarrier(n.header)
	}, otel.GetTextMapPropagator())
}
