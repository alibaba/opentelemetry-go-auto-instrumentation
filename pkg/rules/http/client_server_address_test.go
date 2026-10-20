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

package http

import (
	"net/url"
	"testing"
)

func TestNetHttpClientGetServerAddress(t *testing.T) {
	getter := netHttpClientAttrsGetter{}
	mustURL := func(raw string) *url.URL {
		u, err := url.Parse(raw)
		if err != nil {
			t.Fatalf("parse url: %v", err)
		}
		return u
	}

	tests := []struct {
		name string
		req  *netHttpRequest
		want string
	}{
		{
			name: "host set",
			req:  &netHttpRequest{host: "example.com", url: mustURL("http://other.example/path")},
			want: "example.com",
		},
		{
			name: "host with port preferred over url",
			req:  &netHttpRequest{host: "example.com:8080", url: mustURL("http://other.example/path")},
			want: "example.com",
		},
		{
			name: "empty host falls back to url host",
			req:  &netHttpRequest{url: mustURL("http://example.com/path")},
			want: "example.com",
		},
		{
			name: "empty host falls back to url hostname without port",
			req:  &netHttpRequest{url: mustURL("http://example.com:8080/path")},
			want: "example.com",
		},
		{
			name: "ipv4 host with port",
			req:  &netHttpRequest{host: "192.168.1.10:8080", url: mustURL("http://other.example/path")},
			want: "192.168.1.10",
		},
		{
			name: "empty host falls back to ipv4 url without port",
			req:  &netHttpRequest{url: mustURL("http://10.0.0.1/v1")},
			want: "10.0.0.1",
		},
		{
			name: "empty host falls back to ipv4 url with port",
			req:  &netHttpRequest{url: mustURL("http://192.168.1.10:8080/v1")},
			want: "192.168.1.10",
		},
		{
			name: "ipv6 host with port",
			req:  &netHttpRequest{host: "[2001:db8::1]:8443", url: mustURL("http://other.example/path")},
			want: "2001:db8::1",
		},
		{
			name: "ipv6 host without port",
			req:  &netHttpRequest{host: "[::1]", url: mustURL("http://other.example/path")},
			want: "::1",
		},
		{
			name: "empty host falls back to ipv6 url with port",
			req:  &netHttpRequest{url: mustURL("http://[2001:db8::1]:8443/v1")},
			want: "2001:db8::1",
		},
		{
			name: "empty host falls back to ipv6 url without port",
			req:  &netHttpRequest{url: mustURL("http://[::1]/v1")},
			want: "::1",
		},
		{
			name: "empty host and nil url",
			req:  &netHttpRequest{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getter.GetServerAddress(tt.req); got != tt.want {
				t.Fatalf("GetServerAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}
