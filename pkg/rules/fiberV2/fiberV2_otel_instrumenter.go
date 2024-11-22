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
package fiberV2

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"strconv"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/http"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/net"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var emptyFiberV2Response = fiberV2Response{}

var fiberV2Enabler = instrumenter.NewDefaultInstrumentEnabler()

type fiberV2ServerAttrsGetter struct {
}

func (n fiberV2ServerAttrsGetter) GetRequestMethod(request *fiberV2Request) string {
	return request.method
}
func (n fiberV2ServerAttrsGetter) GetHttpRequestHeader(request *fiberV2Request, name string) []string {
	all := make([]string, 0)
	for _, header := range request.header.PeekAll(name) {
		all = append(all, string(header))
	}
	return all
}
func (n fiberV2ServerAttrsGetter) GetHttpResponseStatusCode(request *fiberV2Request, response *fiberV2Response, err error) int {
	return response.statusCode
}
func (n fiberV2ServerAttrsGetter) GetHttpResponseHeader(request *fiberV2Request, response *fiberV2Response, name string) []string {
	all := make([]string, 0)
	for _, header := range response.header.PeekAll(name) {
		all = append(all, string(header))
	}
	return all
}
func (n fiberV2ServerAttrsGetter) GetErrorType(request *fiberV2Request, response *fiberV2Response, err error) string {
	return ""
}
func (n fiberV2ServerAttrsGetter) GetUrlScheme(request *fiberV2Request) string {
	return request.url.Scheme
}
func (n fiberV2ServerAttrsGetter) GetUrlPath(request *fiberV2Request) string {
	return request.url.Path
}
func (n fiberV2ServerAttrsGetter) GetUrlQuery(request *fiberV2Request) string {
	return request.url.RawQuery
}
func (n fiberV2ServerAttrsGetter) GetNetworkType(request *fiberV2Request, response *fiberV2Response) string {
	return "ipv4"
}
func (n fiberV2ServerAttrsGetter) GetNetworkTransport(request *fiberV2Request, response *fiberV2Response) string {
	return "tcp"
}
func (n fiberV2ServerAttrsGetter) GetNetworkProtocolName(request *fiberV2Request, response *fiberV2Response) string {
	if request.isTls == false {
		return "http"
	} else {
		return "https"
	}
}
func (n fiberV2ServerAttrsGetter) GetNetworkProtocolVersion(request *fiberV2Request, response *fiberV2Response) string {
	return ""
}
func (n fiberV2ServerAttrsGetter) GetNetworkLocalInetAddress(request *fiberV2Request, response *fiberV2Response) string {
	return ""
}
func (n fiberV2ServerAttrsGetter) GetNetworkLocalPort(request *fiberV2Request, response *fiberV2Response) int {
	return 0
}
func (n fiberV2ServerAttrsGetter) GetNetworkPeerInetAddress(request *fiberV2Request, response *fiberV2Response) string {
	return request.url.Host
}
func (n fiberV2ServerAttrsGetter) GetNetworkPeerPort(request *fiberV2Request, response *fiberV2Response) int {
	port, err := strconv.Atoi(request.url.Port())
	if err != nil {
		return 0
	}
	return port
}
func (n fiberV2ServerAttrsGetter) GetHttpRoute(request *fiberV2Request) string {
	return request.url.Path
}

type fiberV2RequestCarrier struct {
	req *fasthttp.RequestHeader
}

func (f fiberV2RequestCarrier) Get(key string) string {
	return string(f.req.Peek(key))
}
func (f fiberV2RequestCarrier) Set(key string, value string) {
	f.req.Set(key, value)
}
func (f fiberV2RequestCarrier) Keys() []string {
	keyStrs := make([]string, 0)
	peekKeys := f.req.PeekKeys()
	for _, peekKey := range peekKeys {
		keyStrs = append(keyStrs, string(peekKey))
	}
	return keyStrs
}

func BuildFiberV2ServerOtelInstrumenter() *instrumenter.PropagatingFromUpstreamInstrumenter[*fiberV2Request, *fiberV2Response] {
	builder := instrumenter.Builder[*fiberV2Request, *fiberV2Response]{}
	serverGetter := fiberV2ServerAttrsGetter{}
	commonExtractor := http.HttpCommonAttrsExtractor[*fiberV2Request, *fiberV2Response, fiberV2ServerAttrsGetter, fiberV2ServerAttrsGetter]{HttpGetter: serverGetter, NetGetter: serverGetter}
	networkExtractor := net.NetworkAttrsExtractor[*fiberV2Request, *fiberV2Response, fiberV2ServerAttrsGetter]{Getter: serverGetter}
	urlExtractor := net.UrlAttrsExtractor[*fiberV2Request, *fiberV2Response, fiberV2ServerAttrsGetter]{Getter: serverGetter}
	return builder.Init().SetSpanStatusExtractor(http.HttpServerSpanStatusExtractor[*fiberV2Request, *fiberV2Response]{Getter: serverGetter}).SetSpanNameExtractor(&http.HttpServerSpanNameExtractor[*fiberV2Request, *fiberV2Response]{Getter: serverGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysServerExtractor[*fiberV2Request]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.FIBER_V2_SERVER_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddAttributesExtractor(&http.HttpServerAttrsExtractor[*fiberV2Request, *fiberV2Response, fiberV2ServerAttrsGetter, fiberV2ServerAttrsGetter, fiberV2ServerAttrsGetter]{Base: commonExtractor, NetworkExtractor: networkExtractor, UrlExtractor: urlExtractor}).BuildPropagatingFromUpstreamInstrumenter(func(n *fiberV2Request) propagation.TextMapCarrier {
		return fiberV2RequestCarrier{req: n.header}
	}, otel.GetTextMapPropagator())
}
