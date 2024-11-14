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

package mux

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	mux "github.com/gorilla/mux"
)

var muxEnabler = instrumenter.NewDefaultInstrumentEnabler()

func muxRoute130OnEnter(call api.CallContext, req *http.Request, route interface{}) {
	if !muxEnabler.Enable() {
		return
	}
	if req != nil {
		lcs := trace.LocalRootSpanFromGLS()
		if lcs != nil && route != nil {
			r, ok := route.(*mux.Route)
			if ok {
				tmpl, err := r.GetPathTemplate()
				if err == nil && req.URL != nil && tmpl != req.URL.Path {
					lcs.SetName(tmpl)
				}
			}
		}
	}
}

// since mux v1.7.4
func muxRoute174OnEnter(call api.CallContext, req *http.Request, route *mux.Route) {
	if !muxEnabler.Enable() {
		return
	}
	if req != nil {
		lcs := trace.LocalRootSpanFromGLS()
		if lcs != nil && route != nil {
			tmpl, err := route.GetPathTemplate()
			if err == nil && req.URL != nil && tmpl != req.URL.Path {
				lcs.SetName(tmpl)
			}
		}
	}
}
