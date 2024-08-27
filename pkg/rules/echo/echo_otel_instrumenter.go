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

type echoServerAttrsGetter struct {
}

func (n echoServerAttrsGetter) GetRequestMethod(request echoRequest) string {
	return request.method
}

func (n echoServerAttrsGetter) GetHttpRequestHeader(request echoRequest, name string) []string {
	return request.header.Values(name)
}

func (n echoServerAttrsGetter) GetHttpResponseStatusCode(request echoRequest, response echoResponse, err error) int {
	return response.statusCode
}

func (n echoServerAttrsGetter) GetHttpResponseHeader(request echoRequest, response echoResponse, name string) []string {
	return response.header.Values(name)
}

func (n echoServerAttrsGetter) GetErrorType(request echoRequest, response echoResponse, err error) string {
	return ""
}

func (n echoServerAttrsGetter) GetUrlScheme(request echoRequest) string {
	return request.url.Scheme
}

func (n echoServerAttrsGetter) GetUrlPath(request echoRequest) string {
	return request.url.Path
}

func (n echoServerAttrsGetter) GetUrlQuery(request echoRequest) string {
	return request.url.RawQuery
}

func (n echoServerAttrsGetter) GetNetworkType(request echoRequest, response echoResponse) string {
	return "ipv4"
}

func (n echoServerAttrsGetter) GetNetworkTransport(request echoRequest, response echoResponse) string {
	return "tcp"
}

func (n echoServerAttrsGetter) GetNetworkProtocolName(request echoRequest, response echoResponse) string {
	if request.isTls == false {
		return "http"
	} else {
		return "https"
	}
}

func (n echoServerAttrsGetter) GetNetworkProtocolVersion(request echoRequest, response echoResponse) string {
	return request.version
}

func (n echoServerAttrsGetter) GetNetworkLocalInetAddress(request echoRequest, response echoResponse) string {
	return ""
}

func (n echoServerAttrsGetter) GetNetworkLocalPort(request echoRequest, response echoResponse) int {
	return 0
}

func (n echoServerAttrsGetter) GetNetworkPeerInetAddress(request echoRequest, response echoResponse) string {
	return request.host
}

func (n echoServerAttrsGetter) GetNetworkPeerPort(request echoRequest, response echoResponse) int {
	port, err := strconv.Atoi(request.url.Port())
	if err != nil {
		return 0
	}
	return port
}

func (n echoServerAttrsGetter) GetHttpRoute(request echoRequest) string {
	return request.path
}

func BuildEchoServerOtelInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[echoRequest, echoResponse] {
	builder := instrumenter.Builder[echoRequest, echoResponse]{}
	serverGetter := echoServerAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[echoRequest, echoResponse, echoServerAttrsGetter, echoServerAttrsGetter]{HttpGetter: serverGetter, NetGetter: serverGetter, Converter: &http.ServerHttpStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[echoRequest, echoResponse, echoServerAttrsGetter]{Getter: serverGetter}
	urlExtractor := net.UrlAttrsExtractor[echoRequest, echoResponse, echoServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanNameExtractor(&http.HttpServerSpanNameExtractor[echoRequest, echoResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[echoRequest]{}).
		AddAttributesExtractor(&http.HttpServerAttrsExtractor[echoRequest, echoResponse, echoServerAttrsGetter, echoServerAttrsGetter, echoServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor, UrlExtractor: urlExtractor}).BuildPropagatingFromUpstreamInstrumenter(func(n echoRequest) propagation.TextMapCarrier {
		if n.header == nil {
			return nil
		}
		return propagation.HeaderCarrier(n.header)
	}, otel.GetTextMapPropagator())
}
