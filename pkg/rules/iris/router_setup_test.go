// Copyright (c) 2026 Alibaba Group Holding Ltd.
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
	"testing"

	"github.com/kataras/iris/v12/macro"
)

type stubIrisRoute struct {
	statusErrorCode int
	tmplSrc         string
	path            string
}

func (s stubIrisRoute) StatusErrorCode() int { return s.statusErrorCode }
func (s stubIrisRoute) Tmpl() macro.Template {
	return macro.Template{Src: s.tmplSrc}
}
func (s stubIrisRoute) Path() string { return s.path }

func TestRouteTemplateFromIrisRoute(t *testing.T) {
	tests := []struct {
		name string
		curr irisRouteInfo
		want string
	}{
		{
			name: "nil route",
			curr: nil,
			want: "",
		},
		{
			name: "resource route uses template src",
			curr: stubIrisRoute{tmplSrc: "/users/{id:int}", path: "/users/{id:int}"},
			want: "/users/{id:int}",
		},
		{
			name: "resource route falls back to path when tmpl empty",
			curr: stubIrisRoute{path: "/iris"},
			want: "/iris",
		},
		{
			name: "error handler route is skipped",
			curr: stubIrisRoute{statusErrorCode: 404, tmplSrc: "/", path: "/"},
			want: "",
		},
		{
			name: "error handler does not leak slash template",
			curr: stubIrisRoute{statusErrorCode: 500, tmplSrc: "/", path: "/"},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := routeTemplateFromIrisRoute(tt.curr)
			if got != tt.want {
				t.Fatalf("routeTemplateFromIrisRoute() = %q, want %q", got, tt.want)
			}
		})
	}
}
