// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
