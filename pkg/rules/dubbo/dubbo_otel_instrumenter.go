// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:build ignore

package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/dubbo"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"strconv"
	"strings"
)

type dubboClientAttrsGetter struct {
}

func (n dubboClientAttrsGetter) GetComponentName(request dubboRequest) string {
	return "dubbo_client"
}

func (n dubboClientAttrsGetter) GetRequestMethod(request dubboRequest) string {
	return request.method
}

func (n dubboClientAttrsGetter) GetNetworkType(request dubboRequest, response dubboResponse) string {
	return "ipv4"
}

func (n dubboClientAttrsGetter) GetNetworkTransport(request dubboRequest, response dubboResponse) string {
	return "tcp"
}

func (n dubboClientAttrsGetter) GetNetworkProtocolName(request dubboRequest, response dubboResponse) string {
	return "dubbo"
}

func (n dubboClientAttrsGetter) GetNetworkProtocolVersion(request dubboRequest, response dubboResponse) string {
	return "3.0"
}

func (n dubboClientAttrsGetter) GetNetworkLocalInetAddress(request dubboRequest, response dubboResponse) string {
	return ""
}

func (n dubboClientAttrsGetter) GetNetworkLocalPort(request dubboRequest, response dubboResponse) int {
	return 0
}

func (n dubboClientAttrsGetter) GetNetworkPeerInetAddress(request dubboRequest, response dubboResponse) string {
	return request.addr
}

func (n dubboClientAttrsGetter) GetNetworkPeerPort(request dubboRequest, response dubboResponse) int {
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

func (n dubboClientAttrsGetter) GetServerAddress(request dubboRequest) string {
	return request.addr
}

func (n dubboClientAttrsGetter) GetErrorType(request dubboRequest, response dubboResponse, err error) string {
	return ""
}

func (n dubboClientAttrsGetter) GetDubboResponseStatusCode(request dubboRequest, response dubboResponse, err error) string {
	return response.statusCode
}

type dubboServerAttrsGetter struct {
}

func (n dubboServerAttrsGetter) GetComponentName(request dubboRequest) string {
	return "dubbo_server"
}

func (n dubboServerAttrsGetter) GetRequestMethod(request dubboRequest) string {
	return request.method
}

func (n dubboServerAttrsGetter) GetErrorType(request dubboRequest, response dubboResponse, err error) string {
	return ""
}

func (n dubboServerAttrsGetter) GetDubboResponseStatusCode(request dubboRequest, response dubboResponse, err error) string {
	return response.statusCode
}

func (n dubboServerAttrsGetter) GetNetworkType(request dubboRequest, response dubboResponse) string {
	return "ipv4"
}

func (n dubboServerAttrsGetter) GetNetworkTransport(request dubboRequest, response dubboResponse) string {
	return "tcp"
}

func (n dubboServerAttrsGetter) GetNetworkProtocolName(request dubboRequest, response dubboResponse) string {
	return "dubbo"
}

func (n dubboServerAttrsGetter) GetNetworkProtocolVersion(request dubboRequest, response dubboResponse) string {
	return "3.0"
}

func (n dubboServerAttrsGetter) GetNetworkLocalInetAddress(request dubboRequest, response dubboResponse) string {
	return ""
}

func (n dubboServerAttrsGetter) GetNetworkLocalPort(request dubboRequest, response dubboResponse) int {
	return 0
}

func (n dubboServerAttrsGetter) GetNetworkPeerInetAddress(request dubboRequest, response dubboResponse) string {
	return request.addr
}

func (n dubboServerAttrsGetter) GetNetworkPeerPort(request dubboRequest, response dubboResponse) int {
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

type dubboClientSpanNameExtractor struct {
}

func (n dubboClientSpanNameExtractor) Extract(request dubboRequest) string {
	return request.method
}

type dubboServerSpanNameExtractor struct {
}

func (n dubboServerSpanNameExtractor) Extract(request dubboRequest) string {
	return request.method
}

func BuildDubboClientInstrumenter() *instrumenter.PropagatingToDownstreamInstrumenter[dubboRequest, dubboResponse] {
	builder := instrumenter.Builder[dubboRequest, dubboResponse]{}
	clientGetter := dubboClientAttrsGetter{}
	commonExtractor := dubbo.DubboCommonAttrsExtractor[dubboRequest, dubboResponse, dubboClientAttrsGetter, dubboClientAttrsGetter]{DubboGetter: clientGetter, NetGetter: clientGetter, Converter: &dubbo.ClientDubboStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[dubboRequest, dubboResponse, dubboClientAttrsGetter]{Getter: clientGetter}
	return builder.Init().SetSpanNameExtractor(&dubbo.DubboClientSpanNameExtractor[dubboRequest, dubboResponse]{Getter: clientGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[dubboRequest]{}).
		AddAttributesExtractor(&dubbo.DubboClientAttrsExtractor[dubboRequest, dubboResponse, dubboClientAttrsGetter, dubboClientAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor}).
		BuildPropagatingToDownstreamInstrumenter(func(n dubboRequest) propagation.TextMapCarrier {
			if n.metadata == nil {
				return nil
			}
			return &metadataSupplier{
				metadata: n.metadata,
			}
		}, otel.GetTextMapPropagator())
}

func BuildDubboServerInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[dubboRequest, dubboResponse] {
	builder := instrumenter.Builder[dubboRequest, dubboResponse]{}
	serverGetter := dubboServerAttrsGetter{}
	commonExtractor := dubbo.DubboCommonAttrsExtractor[dubboRequest, dubboResponse, dubboServerAttrsGetter, dubboServerAttrsGetter]{DubboGetter: serverGetter, NetGetter: serverGetter, Converter: &dubbo.ServerDubboStatusCodeConverter{}}
	networkExtractor := net.NetworkAttrsExtractor[dubboRequest, dubboResponse, dubboServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanNameExtractor(&dubbo.DubboServerSpanNameExtractor[dubboRequest, dubboResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[dubboRequest]{}).
		AddAttributesExtractor(&dubbo.DubboServerAttrsExtractor[dubboRequest, dubboResponse, dubboServerAttrsGetter, dubboServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor}).
		BuildPropagatingFromUpstreamInstrumenter(func(n dubboRequest) propagation.TextMapCarrier {
			if n.metadata == nil {
				return nil
			}
			return &metadataSupplier{
				metadata: n.metadata,
			}
		}, otel.GetTextMapPropagator())
}
