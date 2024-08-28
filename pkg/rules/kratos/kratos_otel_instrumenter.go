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
	"go.opentelemetry.io/otel/propagation"
	"strconv"
	"strings"
)

type kratosClientAttrsGetter struct {
}

func (n kratosClientAttrsGetter) GetComponentName(request kratosRequest) string {
	return request.componentName
}

func (n kratosClientAttrsGetter) GetErrorType(request kratosRequest, response kratosResponse, err error) string {
	return ""
}

func (n kratosClientAttrsGetter) GetNetworkType(request kratosRequest, response kratosResponse) string {
	return "ipv4"
}

func (n kratosClientAttrsGetter) GetNetworkTransport(request kratosRequest, response kratosResponse) string {
	return "tcp"
}

func (n kratosClientAttrsGetter) GetNetworkProtocolName(request kratosRequest, response kratosResponse) string {
	if request.httpMethod != "" {
		return "http"
	} else {
		return "grpc"
	}
}

func (n kratosClientAttrsGetter) GetNetworkProtocolVersion(request kratosRequest, response kratosResponse) string {
	return ""
}

func (n kratosClientAttrsGetter) GetNetworkLocalInetAddress(request kratosRequest, response kratosResponse) string {
	return ""
}

func (n kratosClientAttrsGetter) GetNetworkLocalPort(request kratosRequest, response kratosResponse) int {
	return 0
}

func (n kratosClientAttrsGetter) GetNetworkPeerInetAddress(request kratosRequest, response kratosResponse) string {
	return request.addr
}

func (n kratosClientAttrsGetter) GetNetworkPeerPort(request kratosRequest, response kratosResponse) int {
	ports := strings.Split(request.addr, ":")
	if len(ports) != 2 {
		return 0
	}
	port, err := strconv.Atoi(ports[1])
	if err != nil {
		return 0
	}
	return port
}

func (n kratosClientAttrsGetter) GetRequestMethod(request kratosRequest) string {
	return request.httpMethod
}

func (n kratosClientAttrsGetter) GetSpanName(request kratosRequest) string {
	return request.method
}

func (n kratosClientAttrsGetter) GetHttpRequestHeader(request kratosRequest, name string) []string {
	return []string{}
}

func (n kratosClientAttrsGetter) GetHttpResponseHeader(request kratosRequest, response kratosResponse, name string) []string {
	return []string{}
}

func (n kratosClientAttrsGetter) GetHttpResponseStatusCode(request kratosRequest, response kratosResponse, err error) int {
	return response.statusCode
}

func (n kratosClientAttrsGetter) GetHttpMethod(request kratosRequest) string {
	return request.httpMethod
}

func (n kratosClientAttrsGetter) GetUrlFull(request kratosRequest) string {
	return ""
}

func (n kratosClientAttrsGetter) GetServerAddress(request kratosRequest) string {
	return request.addr
}

type kratosServerAttrsGetter struct {
}

func (n kratosServerAttrsGetter) GetComponentName(request kratosRequest) string {
	return request.componentName
}

func (n kratosServerAttrsGetter) GetSpanName(request kratosRequest) string {
	return request.method
}

func (n kratosServerAttrsGetter) GetHttpResponseStatusCode(request kratosRequest, response kratosResponse, err error) int {
	return response.statusCode
}

func (n kratosServerAttrsGetter) GetHttpRequestHeader(request kratosRequest, name string) []string {
	return []string{}
}

func (n kratosServerAttrsGetter) GetHttpResponseHeader(request kratosRequest, response kratosResponse, name string) []string {
	return []string{}
}

func (n kratosServerAttrsGetter) GetHttpMethod(request kratosRequest) string {
	return request.httpMethod
}

func (n kratosServerAttrsGetter) GetErrorType(request kratosRequest, response kratosResponse, err error) string {
	return ""
}

func (n kratosServerAttrsGetter) GetNetworkType(request kratosRequest, response kratosResponse) string {
	return "ipv4"
}

func (n kratosServerAttrsGetter) GetNetworkTransport(request kratosRequest, response kratosResponse) string {
	return "tcp"
}

func (n kratosServerAttrsGetter) GetNetworkProtocolName(request kratosRequest, response kratosResponse) string {
	if request.httpMethod != "" {
		return "http"
	} else {
		return "grpc"
	}
}

func (n kratosServerAttrsGetter) GetNetworkProtocolVersion(request kratosRequest, response kratosResponse) string {
	return ""
}

func (n kratosServerAttrsGetter) GetNetworkLocalInetAddress(request kratosRequest, response kratosResponse) string {
	return ""
}

func (n kratosServerAttrsGetter) GetNetworkLocalPort(request kratosRequest, response kratosResponse) int {
	return 0
}

func (n kratosServerAttrsGetter) GetRequestMethod(request kratosRequest) string {
	return request.httpMethod
}

func (n kratosServerAttrsGetter) GetNetworkPeerInetAddress(request kratosRequest, response kratosResponse) string {
	return request.addr
}

func (n kratosServerAttrsGetter) GetUrlScheme(request kratosRequest) string {
	return ""
}

func (n kratosServerAttrsGetter) GetUrlPath(request kratosRequest) string {
	return request.method
}

func (n kratosServerAttrsGetter) GetHttpRoute(request kratosRequest) string {
	return ""
}

func (n kratosServerAttrsGetter) GetUrlQuery(request kratosRequest) string {
	return ""
}

func (n kratosServerAttrsGetter) GetNetworkPeerPort(request kratosRequest, response kratosResponse) int {
	ports := strings.Split(request.addr, ":")
	if len(ports) != 2 {
		return 0
	}
	port, err := strconv.Atoi(ports[1])
	if err != nil {
		return 0
	}
	return port
}

type kratosClientSpanNameExtractor struct {
}

func (n kratosClientSpanNameExtractor) Extract(request kratosRequest) string {
	return request.method
}

type kratosServerSpanNameExtractor struct {
}

func (n kratosServerSpanNameExtractor) Extract(request kratosRequest) string {
	return request.method
}

func BuildKratosClientInstrumenter() *instrumenter.PropagatingToDownstreamInstrumenter[kratosRequest, kratosResponse] {
	builder := instrumenter.Builder[kratosRequest, kratosResponse]{}
	clientGetter := kratosClientAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[kratosRequest, kratosResponse, kratosClientAttrsGetter, kratosClientAttrsGetter]{HttpGetter: clientGetter, NetGetter: clientGetter, Converter: &http.ClientHttpStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[kratosRequest, kratosResponse, kratosClientAttrsGetter]{Getter: clientGetter}
	return builder.Init().SetSpanNameExtractor(&http.HttpClientSpanNameExtractor[kratosRequest, kratosResponse]{Getter: clientGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[kratosRequest]{}).
		AddAttributesExtractor(&http.HttpClientAttrsExtractor[kratosRequest, kratosResponse, kratosClientAttrsGetter, kratosClientAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor}).
		BuildPropagatingToDownstreamInstrumenter(func(n kratosRequest) propagation.TextMapCarrier {
			return n.header
		}, kratosPropagators)
}

func BuildKratosServerInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[kratosRequest, kratosResponse] {
	builder := instrumenter.Builder[kratosRequest, kratosResponse]{}
	serverGetter := kratosServerAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[kratosRequest, kratosResponse, kratosServerAttrsGetter, kratosServerAttrsGetter]{HttpGetter: serverGetter, NetGetter: serverGetter, Converter: &http.ServerHttpStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[kratosRequest, kratosResponse, kratosServerAttrsGetter]{Getter: serverGetter}
	urlExtractor := net.UrlAttrsExtractor[kratosRequest, kratosResponse, kratosServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanNameExtractor(&http.HttpServerSpanNameExtractor[kratosRequest, kratosResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[kratosRequest]{}).
		AddAttributesExtractor(&http.HttpServerAttrsExtractor[kratosRequest, kratosResponse, kratosServerAttrsGetter, kratosServerAttrsGetter, kratosServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor, UrlExtractor: urlExtractor}).
		BuildPropagatingFromUpstreamInstrumenter(func(n kratosRequest) propagation.TextMapCarrier {
			return n.header
		}, kratosPropagators)
}
