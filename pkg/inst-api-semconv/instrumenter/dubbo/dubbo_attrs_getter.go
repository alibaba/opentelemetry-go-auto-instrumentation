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
package dubbo

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"

type DubboCommonAttrsGetter[REQUEST any, RESPONSE any] interface {
	GetDubboResponseStatusCode(request REQUEST, response RESPONSE, err error) string
	GetErrorType(request REQUEST, response RESPONSE, err error) string
	GetRequestMethod(request REQUEST) string
}

type DubboServerAttrsGetter[REQUEST any, RESPONSE any] interface {
	DubboCommonAttrsGetter[REQUEST, RESPONSE]
	net.NetworkAttrsGetter[REQUEST, RESPONSE]
	GetRequestMethod(request REQUEST) string
}

type DubboClientAttrsGetter[REQUEST any, RESPONSE any] interface {
	DubboCommonAttrsGetter[REQUEST, RESPONSE]
	net.NetworkAttrsGetter[REQUEST, RESPONSE]
	GetServerAddress(request REQUEST) string
	GetRequestMethod(request REQUEST) string
}
