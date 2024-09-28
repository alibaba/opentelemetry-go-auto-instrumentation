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
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/http"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"strconv"
)

type ginServerAttrsGetter struct {
}

func (n ginServerAttrsGetter) GetRequestMethod(request ginRequest) string {
	return request.method
}

func (n ginServerAttrsGetter) GetHttpRequestHeader(request ginRequest, name string) []string {
	return request.header.Values(name)
}

func (n ginServerAttrsGetter) GetHttpResponseStatusCode(request ginRequest, response ginResponse, err error) int {
	return response.statusCode
}

func (n ginServerAttrsGetter) GetHttpResponseHeader(request ginRequest, response ginResponse, name string) []string {
	return response.header.Values(name)
}

func (n ginServerAttrsGetter) GetErrorType(request ginRequest, response ginResponse, err error) string {
	return ""
}

func (n ginServerAttrsGetter) GetUrlScheme(request ginRequest) string {
	return request.url.Scheme
}

func (n ginServerAttrsGetter) GetUrlPath(request ginRequest) string {
	return request.url.Path
}

func (n ginServerAttrsGetter) GetUrlQuery(request ginRequest) string {
	return request.url.RawQuery
}

func (n ginServerAttrsGetter) GetNetworkType(request ginRequest, response ginResponse) string {
	return "ipv4"
}

func (n ginServerAttrsGetter) GetNetworkTransport(request ginRequest, response ginResponse) string {
	return "tcp"
}

func (n ginServerAttrsGetter) GetNetworkProtocolName(request ginRequest, response ginResponse) string {
	if request.isTls == false {
		return "http"
	} else {
		return "https"
	}
}

func (n ginServerAttrsGetter) GetNetworkProtocolVersion(request ginRequest, response ginResponse) string {
	return request.version
}

func (n ginServerAttrsGetter) GetNetworkLocalInetAddress(request ginRequest, response ginResponse) string {
	return ""
}

func (n ginServerAttrsGetter) GetNetworkLocalPort(request ginRequest, response ginResponse) int {
	return 0
}

func (n ginServerAttrsGetter) GetNetworkPeerInetAddress(request ginRequest, response ginResponse) string {
	return request.host
}

func (n ginServerAttrsGetter) GetNetworkPeerPort(request ginRequest, response ginResponse) int {
	port, err := strconv.Atoi(request.url.Port())
	if err != nil {
		return 0
	}
	return port
}

func (n ginServerAttrsGetter) GetHttpRoute(request ginRequest) string {
	return request.url.Path
}

func BuildGinServerOtelInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[ginRequest, ginResponse] {
	builder := instrumenter.Builder[ginRequest, ginResponse]{}
	serverGetter := ginServerAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[ginRequest, ginResponse, ginServerAttrsGetter, ginServerAttrsGetter]{HttpGetter: serverGetter, NetGetter: serverGetter}
	networkExtractor := net.NetworkAttrsExtractor[ginRequest, ginResponse, ginServerAttrsGetter]{Getter: serverGetter}
	urlExtractor := net.UrlAttrsExtractor[ginRequest, ginResponse, ginServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanStatusExtractor(http.HttpServerSpanStatusExtractor[ginRequest, ginResponse]{Getter: serverGetter}).SetSpanNameExtractor(&http.HttpServerSpanNameExtractor[ginRequest, ginResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[ginRequest]{}).
		AddAttributesExtractor(&http.HttpServerAttrsExtractor[ginRequest, ginResponse, ginServerAttrsGetter, ginServerAttrsGetter, ginServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor, UrlExtractor: urlExtractor}).BuildPropagatingFromUpstreamInstrumenter(func(n ginRequest) propagation.TextMapCarrier {
		if n.header == nil {
			return nil
		}
		return propagation.HeaderCarrier(n.header)
	}, otel.GetTextMapPropagator())
}
