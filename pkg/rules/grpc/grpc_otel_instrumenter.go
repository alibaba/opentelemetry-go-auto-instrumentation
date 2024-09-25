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
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type grpcClientAttrsGetter struct {
}

func (n grpcClientAttrsGetter) GetNetworkType(request grpcRequest, response grpcResponse) string {
	return "ipv4"
}

func (n grpcClientAttrsGetter) GetNetworkLocalInetAddress(request grpcRequest, response grpcResponse) string {
	return ""
}

func (n grpcClientAttrsGetter) GetNetworkLocalPort(request grpcRequest, response grpcResponse) int {
	return 0
}

func (n grpcClientAttrsGetter) GetNetworkPeerPort(request grpcRequest, response grpcResponse) int {
	return 0
}

func (n grpcClientAttrsGetter) GetNetworkTransport(request grpcRequest, response grpcResponse) string {
	return "tcp"
}

func (n grpcClientAttrsGetter) GetNetworkProtocolName(request grpcRequest, response grpcResponse) string {
	return "GRPC"
}

func (n grpcClientAttrsGetter) GetNetworkProtocolVersion(request grpcRequest, response grpcResponse) string {
	return "HTTP2"
}

func (n grpcClientAttrsGetter) GetUrlFull(request grpcRequest) string {
	return request.methodName
}

func (n grpcClientAttrsGetter) GetNetworkPeerInetAddress(request grpcRequest, response grpcResponse) string {
	return request.addr
}

func (n grpcClientAttrsGetter) GetComponentName(request grpcRequest) string {
	return "grpc_client"
}

func (n grpcClientAttrsGetter) GetRequestMethod(request grpcRequest) string {
	return request.methodName
}

func (n grpcClientAttrsGetter) GetGrpcResponseStatusCode(request grpcRequest, response grpcResponse, err error) int {
	return response.statusCode
}

func (n grpcClientAttrsGetter) GetGrpcMethod(request grpcRequest) string {
	return ""
}

type grpcClientSpanNameExtractor struct {
}

func (g grpcClientSpanNameExtractor) Extract(request grpcRequest) string {
	return request.methodName
}

type grpcServerAttrsGetter struct {
}

func (n grpcServerAttrsGetter) GetComponentName(request grpcRequest) string {
	return "grpc_server"
}

func (n grpcServerAttrsGetter) GetNetworkType(request grpcRequest, response grpcResponse) string {
	return "ipv4"
}

func (n grpcServerAttrsGetter) GetNetworkPeerPort(request grpcRequest, response grpcResponse) int {
	return 0
}

func (n grpcServerAttrsGetter) GetNetworkTransport(request grpcRequest, response grpcResponse) string {
	return "tcp"
}

func (n grpcServerAttrsGetter) GetNetworkPeerInetAddress(request grpcRequest, response grpcResponse) string {
	return request.addr
}

func (n grpcServerAttrsGetter) GetNetworkLocalPort(request grpcRequest, response grpcResponse) int {
	return 0
}

func (n grpcServerAttrsGetter) GetNetworkProtocolName(request grpcRequest, response grpcResponse) string {
	return "GRPC"
}

func (n grpcServerAttrsGetter) GetNetworkLocalInetAddress(request grpcRequest, response grpcResponse) string {
	return ""
}

func (n grpcServerAttrsGetter) GetNetworkProtocolVersion(request grpcRequest, response grpcResponse) string {
	return "HTTP2"
}

func (n grpcServerAttrsGetter) GetRequestMethod(request grpcRequest) string {
	return request.methodName
}

func (n grpcServerAttrsGetter) GetUrlPath(request grpcRequest) string {
	return request.methodName
}

func (n grpcServerAttrsGetter) GetGrpcResponseStatusCode(request grpcRequest, response grpcResponse, err error) int {
	return response.statusCode
}

func (n grpcServerAttrsGetter) GetGrpcMethod(request grpcRequest) string {
	return ""
}

type grpcServerSpanNameExtractor struct {
}

func (n grpcServerSpanNameExtractor) Extract(request grpcRequest) string {
	return request.methodName
}

func BuildGrpcClientInstrumenter() *instrumenter.PropagatingToDownstreamInstrumenter[grpcRequest, grpcResponse] {
	builder := instrumenter.Builder[grpcRequest, grpcResponse]{}
	clientGetter := grpcClientAttrsGetter{}
	networkExtractor := net.NetworkAttrsExtractor[grpcRequest, grpcResponse, grpcClientAttrsGetter]{Getter: clientGetter}
	commonExtractor := GrpcCommonAttrsExtractor[grpcRequest, grpcResponse, grpcClientAttrsGetter, grpcClientAttrsGetter]{GrpcGetter: clientGetter, Converter: &ClientGrpcStatusCodeConverter{}}
	return builder.Init().SetSpanNameExtractor(&GrpcClientSpanNameExtractor[grpcRequest, grpcResponse]{Getter: clientGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[grpcRequest]{}).
		AddAttributesExtractor(&GrpcClientAttrsExtractor[grpcRequest, grpcResponse, grpcClientAttrsGetter, grpcClientAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor}).
		BuildPropagatingToDownstreamInstrumenter(nil, otel.GetTextMapPropagator())
}

func BuildGrpcServerInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[grpcRequest, grpcResponse] {
	builder := instrumenter.Builder[grpcRequest, grpcResponse]{}
	serverGetter := grpcServerAttrsGetter{}
	networkExtractor := net.NetworkAttrsExtractor[grpcRequest, grpcResponse, grpcServerAttrsGetter]{Getter: serverGetter}
	commonExtractor := GrpcCommonAttrsExtractor[grpcRequest, grpcResponse, grpcServerAttrsGetter, grpcServerAttrsGetter]{GrpcGetter: serverGetter, Converter: &ServerGrpcStatusCodeConverter{}}
	return builder.Init().SetSpanNameExtractor(&GrpcServerSpanNameExtractor[grpcRequest, grpcResponse]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[grpcRequest]{}).
		AddAttributesExtractor(&GrpcServerAttrsExtractor[grpcRequest, grpcResponse, grpcServerAttrsGetter, grpcServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor}).
		BuildPropagatingFromUpstreamInstrumenter(func(n grpcRequest) propagation.TextMapCarrier {
			if n.propagators == nil {
				return nil
			}
			return n.propagators
		}, otel.GetTextMapPropagator())
}
