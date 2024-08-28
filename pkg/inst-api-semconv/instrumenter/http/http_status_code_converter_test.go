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

import "testing"

func TestClientHttpStatusCodeConverter(t *testing.T) {
	c := ClientHttpStatusCodeConverter{}
	if c.IsError(200) {
		t.Fatalf("200 should not be an error")
	}
	if !c.IsError(600) || !c.IsError(90) {
		t.Fatalf("600 and 90 should be an error")
	}
}

func TestServerHttpStatusCodeConverter(t *testing.T) {
	c := ServerHttpStatusCodeConverter{}
	if c.IsError(200) {
		t.Fatalf("200 should not be an error")
	}
	if !c.IsError(500) || !c.IsError(90) {
		t.Fatalf("600 and 90 should be an error")
	}
}
