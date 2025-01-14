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

package verifier

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// VerifyDbAttributes TODO: make attribute name to semconv attribute
func VerifyDbAttributes(span tracetest.SpanStub, name, system, address, statement, operation string) {
	NewSpanVerifier().
		HasSpanKind(trace.SpanKindClient).
		HasName(name).
		HasStringAttribute("db.system", system).
		HasStringAttributeContains("server.address", address).
		HasStringAttributeContains("db.query.text", statement).
		HasStringAttribute("db.operation.name", operation).
		Verify(span)
}

func VerifyHttpClientAttributes(span tracetest.SpanStub, name, method, fullUrl, protocolName, protocolVersion, networkTransport, networkType, localAddr, peerAddr string, statusCode, localPort, peerPort int64) {
	NewSpanVerifier().
		HasSpanKind(trace.SpanKindClient).
		HasName(name).
		HasStringAttribute("http.request.method", method).
		HasStringAttribute("url.full", fullUrl).
		HasStringAttribute("network.protocol.name", protocolName).
		HasStringAttribute("network.protocol.version", protocolVersion).
		HasStringAttribute("network.transport", networkTransport).
		HasStringAttribute("network.type", networkType).
		HasStringAttribute("network.local.address", localAddr).
		HasStringAttribute("network.peer.address", peerAddr).
		HasInt64Attribute("http.response.status_code", statusCode).
		HasInt64Attribute("network.peer.port", peerPort).
		ConditionalVerifier(func() bool {
			return localPort > 0
		}, NewSpanVerifier().HasInt64Attribute("network.local.port", localPort)).
		Verify(span)
}

func VerifyHttpClientMetricsAttributes(attrs []attribute.KeyValue, method, serverAddress, errorType, protocolName, protocolVersion string, serverPort, statusCode int) {
	Assert(GetAttribute(attrs, "http.request.method").AsString() == method, "Except method to be %s, got %s", method, GetAttribute(attrs, "http.request.method").AsString())
	Assert(GetAttribute(attrs, "server.address").AsString() == serverAddress, "Except server.address to be %s, got %s", serverAddress, GetAttribute(attrs, "server.address").AsString())
	Assert(GetAttribute(attrs, "error.type").AsString() == errorType, "Except error.type to be %s, got %s", errorType, GetAttribute(attrs, "error.type").AsString())
	Assert(GetAttribute(attrs, "network.protocol.name").AsString() == protocolName, "Except network.protocol.name to be %s, got %s", protocolName, GetAttribute(attrs, "network.protocol.name").AsString())
	Assert(GetAttribute(attrs, "network.protocol.version").AsString() == protocolVersion, "Except network.protocol.version to be %s, got %s", protocolVersion, GetAttribute(attrs, "network.protocol.version").AsString())
	Assert(GetAttribute(attrs, "server.port").AsInt64() == int64(serverPort), "Except server.port to be %d, got %d", serverPort, GetAttribute(attrs, "server.port").AsInt64())
	Assert(GetAttribute(attrs, "http.response.status_code").AsInt64() == int64(statusCode), "Except status code to be %d, got %d", statusCode, GetAttribute(attrs, "http.response.status_code").AsInt64())
}

func VerifyHttpServerAttributes(span tracetest.SpanStub, name, method, protocolName, networkTransport, networkType, localAddr, peerAddr, agent, scheme, path, query, route string, statusCode int64) {
	NewSpanVerifier().
		HasSpanKind(trace.SpanKindServer).
		HasName(name).
		HasStringAttribute("http.request.method", method).
		HasStringAttribute("network.protocol.name", protocolName).
		HasStringAttribute("network.transport", networkTransport).
		HasStringAttribute("network.type", networkType).
		HasStringAttribute("network.local.address", localAddr).
		HasStringAttribute("network.peer.address", peerAddr).
		HasStringAttribute("user_agent.original", agent).
		HasStringAttribute("url.scheme", scheme).
		HasStringAttribute("url.path", path).
		HasStringAttribute("url.query", query).
		HasStringAttribute("http.route", route).
		HasInt64Attribute("http.response.status_code", statusCode).
		Verify(span)
}

func VerifyHttpServerMetricsAttributes(attrs []attribute.KeyValue, method, httpRoute, errorType, protocolName, protocolVersion, urlScheme string, statusCode int) {
	Assert(GetAttribute(attrs, "http.request.method").AsString() == method, "Except method to be %s, got %s", method, GetAttribute(attrs, "http.request.method").AsString())
	Assert(GetAttribute(attrs, "http.route").AsString() == httpRoute, "Except http.route to be %s, got %s", httpRoute, GetAttribute(attrs, "http.route").AsString())
	Assert(GetAttribute(attrs, "error.type").AsString() == errorType, "Except error.type to be %s, got %s", errorType, GetAttribute(attrs, "error.type").AsString())
	Assert(GetAttribute(attrs, "network.protocol.name").AsString() == protocolName, "Except network.protocol.name to be %s, got %s", protocolName, GetAttribute(attrs, "network.protocol.name").AsString())
	Assert(GetAttribute(attrs, "network.protocol.version").AsString() == protocolVersion, "Except network.protocol.version to be %s, got %s", protocolVersion, GetAttribute(attrs, "network.protocol.version").AsString())
	Assert(GetAttribute(attrs, "url.scheme").AsString() == urlScheme, "Except url.scheme to be %s, got %s", urlScheme, GetAttribute(attrs, "url.scheme").AsString())
	Assert(GetAttribute(attrs, "http.response.status_code").AsInt64() == int64(statusCode), "Except status code to be %d, got %d", statusCode, GetAttribute(attrs, "http.response.status_code").AsInt64())
}

func VerifyRpcServerAttributes(span tracetest.SpanStub, name, system, service, method string) {
	NewSpanVerifier().
		HasSpanKind(trace.SpanKindServer).
		Merge(newRpcSpanVerifier(name, system, service, method)).
		Verify(span)
}

func VerifyRpcClientAttributes(span tracetest.SpanStub, name, system, service, method string) {
	NewSpanVerifier().
		HasSpanKind(trace.SpanKindServer).
		Merge(newRpcSpanVerifier(name, system, service, method)).
		Verify(span)
}

func newRpcSpanVerifier(name, system, service, method string) *SpanVerifier {
	return NewSpanVerifier().
		HasName(name).
		HasStringAttribute("rpc.system", system).
		HasStringAttribute("rpc.service", service).
		HasStringAttribute("rpc.method", method)
}
