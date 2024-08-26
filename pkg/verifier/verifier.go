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
package verifier

import (
	"strings"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// VerifyDbAttributes TODO: make attribute name to semconv attribute
func VerifyDbAttributes(span tracetest.SpanStub, name, dbName, system, user, connString, statement, operation string) {
	Assert(span.SpanKind == trace.SpanKindClient, "Expect to be client span, got %d", span.SpanKind)
	Assert(span.Name == name, "Except client span name to be %s, got %s", name, span.Name)
	actualDbName := GetAttribute(span.Attributes, "db.name").AsString()
	Assert(actualDbName == dbName, "Except client db name to be %s, got %s", dbName, actualDbName)
	actualSystem := GetAttribute(span.Attributes, "db.system").AsString()
	Assert(actualSystem == system, "Except client db system to be %s, got %s", system, actualSystem)
	actualUser := GetAttribute(span.Attributes, "db.user").AsString()
	if actualUser != "" {
		Assert(actualUser == user, "Except client db user to be %s, got %s", user, actualUser)
	}
	actualConnStr := GetAttribute(span.Attributes, "db.connection_string").AsString()
	Assert(strings.Contains(actualConnStr, connString), "Except client db conn str to be %s, got %s", connString, actualConnStr)
	actualStatement := GetAttribute(span.Attributes, "db.statement").AsString()
	Assert(actualStatement == statement, "Except client db statement to be %s, got %s", statement, actualStatement)
	actualOperation := GetAttribute(span.Attributes, "db.operation").AsString()
	Assert(actualOperation == operation, "Except client db operation to be %s, got %s", operation, actualOperation)
}

func VerifyHttpClientAttributes(span tracetest.SpanStub, name, method, fullUrl, protocolName, protocolVersion, networkTransport, networkType, localAddr, peerAddr string, statusCode, localPort, peerPort int64) {
	Assert(span.SpanKind == trace.SpanKindClient, "Expect to be client span, got %d", span.SpanKind)
	Assert(span.Name == name, "Except client span name to be %s, got %s", name, span.Name)
	Assert(GetAttribute(span.Attributes, "http.request.method").AsString() == method, "Except method to be %s, got %s", method, GetAttribute(span.Attributes, "http.request.method").AsString())
	Assert(GetAttribute(span.Attributes, "url.full").AsString() == fullUrl, "Except full url to be %s, got %s", fullUrl, GetAttribute(span.Attributes, "url.full").AsString())
	Assert(GetAttribute(span.Attributes, "network.protocol.name").AsString() == protocolName, "Except protocol name to be %s, got %s", protocolName, GetAttribute(span.Attributes, "network.protocol.name").AsString())
	Assert(GetAttribute(span.Attributes, "network.protocol.version").AsString() == protocolVersion, "Except protocol version to be %s, got %s", protocolVersion, GetAttribute(span.Attributes, "network.protocol.version").AsString())
	Assert(GetAttribute(span.Attributes, "network.transport").AsString() == networkTransport, "Except network transport to be %s, got %s", networkTransport, GetAttribute(span.Attributes, "network.transport").AsString())
	Assert(GetAttribute(span.Attributes, "network.type").AsString() == networkType, "Except network type to be %s, got %s", networkType, GetAttribute(span.Attributes, "network.type").AsString())
	Assert(GetAttribute(span.Attributes, "network.local.address").AsString() == localAddr, "Except local addr to be %s, got %s", localAddr, GetAttribute(span.Attributes, "network.local.address").AsString())
	Assert(GetAttribute(span.Attributes, "network.peer.address").AsString() == peerAddr, "Except peer addr to be %s, got %s", peerAddr, GetAttribute(span.Attributes, "network.peer.address").AsString())
	Assert(GetAttribute(span.Attributes, "http.response.status_code").AsInt64() == statusCode, "Except status code to be %d, got %d", statusCode, GetAttribute(span.Attributes, "http.response.status_code").AsInt64())
	Assert(GetAttribute(span.Attributes, "network.peer.port").AsInt64() == peerPort, "Except peer port to be %d, got %d", peerPort, GetAttribute(span.Attributes, "network.peer.port").AsInt64())
	if localPort > 0 {
		Assert(GetAttribute(span.Attributes, "network.local.port").AsInt64() == localPort, "Except local port to be %d, got %d", localPort, GetAttribute(span.Attributes, "network.local.port").AsInt64())
	}
}

func VerifyHttpServerAttributes(span tracetest.SpanStub, name, method, protocolName, networkTransport, networkType, localAddr, peerAddr, agent, scheme, path, query, route string, statusCode int64) {
	Assert(span.SpanKind == trace.SpanKindServer, "Expect to be client span, got %d", span.SpanKind)
	Assert(span.Name == name, "Except client span name to be %s, got %s", name, span.Name)
	Assert(GetAttribute(span.Attributes, "http.request.method").AsString() == method, "Except method to be %s, got %s", method, GetAttribute(span.Attributes, "http.request.method").AsString())
	Assert(GetAttribute(span.Attributes, "network.protocol.name").AsString() == protocolName, "Except protocol name to be %s, got %s", protocolName, GetAttribute(span.Attributes, "network.protocol.name").AsString())
	Assert(GetAttribute(span.Attributes, "network.transport").AsString() == networkTransport, "Except network transport to be %s, got %s", networkTransport, GetAttribute(span.Attributes, "network.transport").AsString())
	Assert(GetAttribute(span.Attributes, "network.type").AsString() == networkType, "Except network type to be %s, got %s", networkType, GetAttribute(span.Attributes, "network.type").AsString())
	Assert(GetAttribute(span.Attributes, "network.local.address").AsString() == localAddr, "Except local addr to be %s, got %s", localAddr, GetAttribute(span.Attributes, "network.local.address").AsString())
	Assert(GetAttribute(span.Attributes, "network.peer.address").AsString() == peerAddr, "Except peer addr to be %s, got %s", peerAddr, GetAttribute(span.Attributes, "network.peer.address").AsString())
	Assert(GetAttribute(span.Attributes, "user_agent.original").AsString() == agent, "Except user agent to be %s, got %s", agent, GetAttribute(span.Attributes, "user_agent.original").AsString())
	Assert(GetAttribute(span.Attributes, "url.scheme").AsString() == scheme, "Except url scheme to be %s, got %s", scheme, GetAttribute(span.Attributes, "url.scheme").AsString())
	Assert(GetAttribute(span.Attributes, "url.path").AsString() == path, "Except url path to be %s, got %s", path, GetAttribute(span.Attributes, "url.path").AsString())
	Assert(GetAttribute(span.Attributes, "url.query").AsString() == query, "Except url query to be %s, got %s", query, GetAttribute(span.Attributes, "url.query").AsString())
	Assert(GetAttribute(span.Attributes, "http.route").AsString() == route, "Except http route to be %s, got %s", route, GetAttribute(span.Attributes, "http.route").AsString())
	Assert(GetAttribute(span.Attributes, "http.response.status_code").AsInt64() == statusCode, "Except status code to be %d, got %d", statusCode, GetAttribute(span.Attributes, "http.response.status_code").AsInt64())
}

func VerifyGrpcClientAttributes(span tracetest.SpanStub, name, method, protocolName, networkTransport, networkType, localAddr, peerAddr string, statusCode int64) {
	Assert(span.SpanKind == trace.SpanKindClient, "Expect to be client span, got %d", span.SpanKind)
	Assert(span.Name == name, "Except client span name to be %s, got %s", name, span.Name)
	Assert(GetAttribute(span.Attributes, "network.protocol.name").AsString() == protocolName, "Except protocol name to be %s, got %s", protocolName, GetAttribute(span.Attributes, "network.protocol.name").AsString())
	Assert(GetAttribute(span.Attributes, "network.transport").AsString() == networkTransport, "Except network transport to be %s, got %s", networkTransport, GetAttribute(span.Attributes, "network.transport").AsString())
	Assert(GetAttribute(span.Attributes, "network.type").AsString() == networkType, "Except network type to be %s, got %s", networkType, GetAttribute(span.Attributes, "network.type").AsString())
	Assert(GetAttribute(span.Attributes, "network.local.address").AsString() == localAddr, "Except local addr to be %s, got %s", localAddr, GetAttribute(span.Attributes, "network.local.address").AsString())
	Assert(GetAttribute(span.Attributes, "network.peer.address").AsString() == peerAddr, "Except peer addr to be %s, got %s", peerAddr, GetAttribute(span.Attributes, "network.peer.address").AsString())
	Assert(GetAttribute(span.Attributes, "grpc.response.status_code").AsInt64() == statusCode, "Except status code to be %d, got %d", statusCode, GetAttribute(span.Attributes, "grpc.response.status_code").AsInt64())
	Assert(GetAttribute(span.Attributes, "grpc.request.method").AsString() == method, "Except method to be %s, got %s", method, GetAttribute(span.Attributes, "grpc.request.method").AsString())
}

func VerifyGrpcServerAttributes(span tracetest.SpanStub, name, method, protocolName, networkTransport, networkType, localAddr, peerAddr string, statusCode int64) {
	Assert(span.SpanKind == trace.SpanKindServer, "Expect to be grpc server span, got %d", span.SpanKind)
	Assert(span.Name == name, "Except grpc server span name to be %s, got %s", name, span.Name)
	Assert(GetAttribute(span.Attributes, "network.protocol.name").AsString() == protocolName, "Except protocol name to be %s, got %s", protocolName, GetAttribute(span.Attributes, "network.protocol.name").AsString())
	Assert(GetAttribute(span.Attributes, "network.transport").AsString() == networkTransport, "Except network transport to be %s, got %s", networkTransport, GetAttribute(span.Attributes, "network.transport").AsString())
	Assert(GetAttribute(span.Attributes, "network.type").AsString() == networkType, "Except network type to be %s, got %s", networkType, GetAttribute(span.Attributes, "network.type").AsString())
	Assert(GetAttribute(span.Attributes, "network.local.address").AsString() == localAddr, "Except local addr to be %s, got %s", localAddr, GetAttribute(span.Attributes, "network.local.address").AsString())
	Assert(GetAttribute(span.Attributes, "network.peer.address").AsString() == peerAddr, "Except peer addr to be %s, got %s", peerAddr, GetAttribute(span.Attributes, "network.peer.address").AsString())
	Assert(GetAttribute(span.Attributes, "grpc.response.status_code").AsInt64() == statusCode, "Except status code to be %d, got %d", statusCode, GetAttribute(span.Attributes, "grpc.response.status_code").AsInt64())
	Assert(GetAttribute(span.Attributes, "grpc.request.method").AsString() == method, "Except method to be %s, got %s", method, GetAttribute(span.Attributes, "grpc.request.method").AsString())
}
