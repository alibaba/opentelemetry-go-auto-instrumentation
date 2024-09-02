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

package dubbo

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"testing"
)

type dubboServerAttrsGetter struct {
}

type dubboClientAttrsGetter struct {
}

type networkAttrsGetter struct {
}

type urlAttrsGetter struct {
}

func (n networkAttrsGetter) GetNetworkType(request testRequest, response testResponse) string {
	return "network-type"
}

func (n networkAttrsGetter) GetNetworkTransport(request testRequest, response testResponse) string {
	return "network-transport"
}

func (n networkAttrsGetter) GetNetworkProtocolName(request testRequest, response testResponse) string {
	return "network-protocol-name"
}

func (n networkAttrsGetter) GetNetworkProtocolVersion(request testRequest, response testResponse) string {
	return "network-protocol-version"
}

func (n networkAttrsGetter) GetNetworkLocalInetAddress(request testRequest, response testResponse) string {
	return "network-local-inet-address"
}

func (n networkAttrsGetter) GetNetworkLocalPort(request testRequest, response testResponse) int {
	return 8080
}

func (n networkAttrsGetter) GetNetworkPeerInetAddress(request testRequest, response testResponse) string {
	return "network-peer-inet-address"
}

func (n networkAttrsGetter) GetNetworkPeerPort(request testRequest, response testResponse) int {
	return 8080
}

func (h dubboClientAttrsGetter) GetRequestMethod(request testRequest) string {
	return request.Method
}

func (h dubboClientAttrsGetter) GetHttpRequestHeader(request testRequest, name string) []string {
	return []string{"request-header"}
}

func (h dubboClientAttrsGetter) GetHttpResponseStatusCode(request testRequest, response testResponse, err error) int {
	return 200
}

func (h dubboClientAttrsGetter) GetHttpResponseHeader(request testRequest, response testResponse, name string) []string {
	return []string{"response-header"}
}

func (h dubboClientAttrsGetter) GetErrorType(request testRequest, response testResponse, err error) string {
	return ""
}

func (h dubboClientAttrsGetter) GetNetworkType(request testRequest, response testResponse) string {
	return "ipv4"
}

func (h dubboClientAttrsGetter) GetNetworkTransport(request testRequest, response testResponse) string {
	return "TCP"
}

func (h dubboClientAttrsGetter) GetNetworkProtocolName(request testRequest, response testResponse) string {
	return "HTTP"
}

func (h dubboClientAttrsGetter) GetNetworkProtocolVersion(request testRequest, response testResponse) string {
	return "HTTP/1.1"
}

func (h dubboClientAttrsGetter) GetNetworkLocalInetAddress(request testRequest, response testResponse) string {
	return "127.0.0.1"
}

func (h dubboClientAttrsGetter) GetNetworkLocalPort(request testRequest, response testResponse) int {
	return 8080
}

func (h dubboClientAttrsGetter) GetNetworkPeerInetAddress(request testRequest, response testResponse) string {
	return "127.0.0.1"
}

func (h dubboClientAttrsGetter) GetNetworkPeerPort(request testRequest, response testResponse) int {
	return 8080
}

func (h dubboClientAttrsGetter) GetDubboResponseStatusCode(request testRequest, response testResponse, err error) string {
	return response.statusCode
}

func (h dubboClientAttrsGetter) GetUrlFull(request testRequest) string {
	return "url-full"
}

func (h dubboClientAttrsGetter) GetServerAddress(request testRequest) string {
	return "server-address"
}

func (h dubboServerAttrsGetter) GetRequestMethod(request testRequest) string {
	return request.Method
}

func (h dubboServerAttrsGetter) GetDubboResponseStatusCode(request testRequest, response testResponse, err error) string {
	return response.statusCode
}

func (h dubboServerAttrsGetter) GetHttpRequestHeader(request testRequest, name string) []string {
	return []string{"request-header"}
}

func (h dubboServerAttrsGetter) GetHttpResponseStatusCode(request testRequest, response testResponse, err error) int {
	return 200
}

func (h dubboServerAttrsGetter) GetHttpResponseHeader(request testRequest, response testResponse, name string) []string {
	return []string{"response-header"}
}

func (h dubboServerAttrsGetter) GetErrorType(request testRequest, response testResponse, err error) string {
	return "error-type"
}
func (h dubboServerAttrsGetter) GetNetworkType(request testRequest, response testResponse) string {
	return "network-type"
}

func (h dubboServerAttrsGetter) GetNetworkTransport(request testRequest, response testResponse) string {
	return "network-transport"
}

func (h dubboServerAttrsGetter) GetNetworkProtocolName(request testRequest, response testResponse) string {
	return "network-protocol-name"
}

func (h dubboServerAttrsGetter) GetNetworkProtocolVersion(request testRequest, response testResponse) string {
	return "network-protocol-version"
}

func (h dubboServerAttrsGetter) GetNetworkLocalInetAddress(request testRequest, response testResponse) string {
	return "127.0.0.1"
}

func (h dubboServerAttrsGetter) GetNetworkLocalPort(request testRequest, response testResponse) int {
	return 8080
}

func (h dubboServerAttrsGetter) GetNetworkPeerInetAddress(request testRequest, response testResponse) string {
	return "127.0.0.1"
}

func (h dubboServerAttrsGetter) GetNetworkPeerPort(request testRequest, response testResponse) int {
	return 8080
}

func TestDubboClientExtractorEnd(t *testing.T) {
	dubboClientExtractor := DubboClientAttrsExtractor[testRequest, testResponse, dubboClientAttrsGetter, networkAttrsGetter]{
		Base:             DubboCommonAttrsExtractor[testRequest, testResponse, dubboClientAttrsGetter, networkAttrsGetter]{},
		NetworkExtractor: net.NetworkAttrsExtractor[testRequest, testResponse, networkAttrsGetter]{},
	}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs = dubboClientExtractor.OnEnd(attrs, parentContext, testRequest{}, testResponse{
		statusCode: "200",
	}, nil)
	if attrs[0].Key != "dubbo.response.status_code" || attrs[0].Value.AsInt64() != 200 {
		t.Fatalf("status code should be 200")
	}
	if attrs[1].Key != semconv.NetworkProtocolNameKey || attrs[1].Value.AsString() != "network-protocol-name" {
		t.Fatalf("wrong network protocol name")
	}
	if attrs[2].Key != semconv.NetworkProtocolVersionKey || attrs[2].Value.AsString() != "network-protocol-version" {
		t.Fatalf("wrong network protocol version")
	}
	if attrs[3].Key != semconv.NetworkTransportKey || attrs[3].Value.AsString() != "network-transport" {
		t.Fatalf("wrong network transport")
	}
	if attrs[4].Key != semconv.NetworkTypeKey || attrs[4].Value.AsString() != "network-type" {
		t.Fatalf("wrong network type")
	}
	if attrs[5].Key != semconv.NetworkProtocolNameKey || attrs[5].Value.AsString() != "network-protocol-name" {
		t.Fatalf("wrong network protocol name")
	}
	if attrs[6].Key != semconv.NetworkProtocolVersionKey || attrs[6].Value.AsString() != "network-protocol-version" {
		t.Fatalf("wrong network protocol version")
	}
	if attrs[7].Key != semconv.NetworkLocalAddressKey || attrs[7].Value.AsString() != "network-local-inet-address" {
		t.Fatalf("wrong network protocol inet address")
	}
	if attrs[8].Key != semconv.NetworkPeerAddressKey || attrs[8].Value.AsString() != "network-peer-inet-address" {
		t.Fatalf("wrong network peer address")
	}
	if attrs[9].Key != semconv.NetworkLocalPortKey || attrs[9].Value.AsInt64() != 8080 {
		t.Fatalf("wrong network local port")
	}
	if attrs[10].Key != semconv.NetworkPeerPortKey || attrs[10].Value.AsInt64() != 8080 {
		t.Fatalf("wrong network peer port")
	}
}

func TestDubboServerExtractorStart(t *testing.T) {
	dubboServerExtractor := DubboServerAttrsExtractor[testRequest, testResponse, dubboServerAttrsGetter, networkAttrsGetter]{
		Base:             DubboCommonAttrsExtractor[testRequest, testResponse, dubboServerAttrsGetter, networkAttrsGetter]{},
		NetworkExtractor: net.NetworkAttrsExtractor[testRequest, testResponse, networkAttrsGetter]{},
	}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs = dubboServerExtractor.OnStart(attrs, parentContext, testRequest{
		Method: "/org.apache.dubbogo.samples.api.Greeter/SayHello",
	})
	if attrs[0].Key != "dubbo.request.method" || attrs[0].Value.AsString() != "/org.apache.dubbogo.samples.api.Greeter/SayHello" {
		t.Fatalf("dubbo method should be /org.apache.dubbogo.samples.api.Greeter/SayHello")
	}
}

func TestDubboServerExtractorEnd(t *testing.T) {
	dubboServerExtractor := DubboServerAttrsExtractor[testRequest, testResponse, dubboServerAttrsGetter, networkAttrsGetter]{
		Base:             DubboCommonAttrsExtractor[testRequest, testResponse, dubboServerAttrsGetter, networkAttrsGetter]{},
		NetworkExtractor: net.NetworkAttrsExtractor[testRequest, testResponse, networkAttrsGetter]{},
	}
	attrs := make([]attribute.KeyValue, 0)
	parentContext := context.Background()
	attrs = dubboServerExtractor.OnEnd(attrs, parentContext, testRequest{}, testResponse{
		statusCode: "200",
	}, nil)
	if attrs[0].Key != "dubbo.response.status_code" || attrs[0].Value.AsInt64() != 200 {
		t.Fatalf("status code should be 200")
	}
	if attrs[1].Key != semconv.NetworkProtocolNameKey || attrs[1].Value.AsString() != "network-protocol-name" {
		t.Fatalf("wrong network protocol name")
	}
	if attrs[2].Key != semconv.NetworkProtocolVersionKey || attrs[2].Value.AsString() != "network-protocol-version" {
		t.Fatalf("wrong network protocol version")
	}
	if attrs[3].Key != semconv.NetworkTransportKey || attrs[3].Value.AsString() != "network-transport" {
		t.Fatalf("wrong network transport")
	}
	if attrs[4].Key != semconv.NetworkTypeKey || attrs[4].Value.AsString() != "network-type" {
		t.Fatalf("wrong network type")
	}
	if attrs[5].Key != semconv.NetworkProtocolNameKey || attrs[5].Value.AsString() != "network-protocol-name" {
		t.Fatalf("wrong network protocol name")
	}
	if attrs[6].Key != semconv.NetworkProtocolVersionKey || attrs[6].Value.AsString() != "network-protocol-version" {
		t.Fatalf("wrong network protocol version")
	}
	if attrs[7].Key != semconv.NetworkLocalAddressKey || attrs[7].Value.AsString() != "network-local-inet-address" {
		t.Fatalf("wrong network protocol inet address")
	}
	if attrs[8].Key != semconv.NetworkPeerAddressKey || attrs[8].Value.AsString() != "network-peer-inet-address" {
		t.Fatalf("wrong network peer address")
	}
	if attrs[9].Key != semconv.NetworkLocalPortKey || attrs[9].Value.AsInt64() != 8080 {
		t.Fatalf("wrong network local port")
	}
	if attrs[10].Key != semconv.NetworkPeerPortKey || attrs[10].Value.AsInt64() != 8080 {
		t.Fatalf("wrong network peer port")
	}
}
