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
package grpc

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"

type GrpcCommonAttrsGetter[REQUEST any, RESPONSE any] interface {
	GetRequestMethod(request REQUEST) string
	GetGrpcResponseStatusCode(request REQUEST, response RESPONSE, err error) int
}

type GrpcServerAttrsGetter[REQUEST any, RESPONSE any] interface {
	GrpcCommonAttrsGetter[REQUEST, RESPONSE]
	GetUrlPath(request REQUEST) string
}

type GrpcClientAttrsGetter[REQUEST any, RESPONSE any] interface {
	GrpcCommonAttrsGetter[REQUEST, RESPONSE]
	net.NetworkAttrsGetter[REQUEST, RESPONSE]
	GetUrlFull(request REQUEST) string
}