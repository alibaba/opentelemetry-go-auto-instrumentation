// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package iris

import (
	_ "unsafe"

	"github.com/alibaba/loongsuite-go/pkg/api"
	otelhttp "github.com/alibaba/loongsuite-go/pkg/rules/http"
	iContext "github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/macro"
)

//go:linkname irisSetCurrentRouteOnEnter github.com/kataras/iris/v12/context.irisSetCurrentRouteOnEnter
func irisSetCurrentRouteOnEnter(call api.CallContext, iCtx *iContext.Context, curr iContext.RouteReadOnly) {
	if !irisEnabler.Enable() || iCtx == nil || curr == nil {
		return
	}
	r := iCtx.Request()
	if r == nil {
		return
	}
	route := routeTemplateFromIrisRoute(curr)
	if route == "" {
		return
	}
	otelhttp.SetServerRouteTemplate(r, route)
	otelhttp.UpdateServerSpanName(r.Method, route)
}

// irisRouteInfo is the RouteReadOnly surface needed to resolve a low-cardinality
// http.route template without depending on request Path().
type irisRouteInfo interface {
	StatusErrorCode() int
	Tmpl() macro.Template
	Path() string
}

// routeTemplateFromIrisRoute returns a route template for selected resource
// routes only. Error handlers (StatusErrorCode != 0) are skipped so they do not
// overwrite a previously matched template with "/".
func routeTemplateFromIrisRoute(curr irisRouteInfo) string {
	if curr == nil || curr.StatusErrorCode() != 0 {
		return ""
	}
	if tmpl := curr.Tmpl(); tmpl.Src != "" {
		return tmpl.Src
	}
	return curr.Path()
}
