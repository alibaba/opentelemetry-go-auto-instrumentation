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

package http

type HttpClientSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	Getter HttpClientAttrsGetter[REQUEST, RESPONSE]
}

func (h *HttpClientSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	method := h.Getter.GetSpanName(request)
	if method == "" {
		return "HTTP"
	}
	return method
}

type HttpServerSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	Getter HttpServerAttrsGetter[REQUEST, RESPONSE]
}

func (h *HttpServerSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	method := h.Getter.GetSpanName(request)
	route := h.Getter.GetHttpRoute(request)
	if method == "" {
		return "HTTP"
	}
	if route == "" {
		return method
	}
	return method + " " + route
}
