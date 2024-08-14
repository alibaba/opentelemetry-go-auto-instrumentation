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
package rpc

type RpcSpanNameExtractor[REQUEST any] struct {
	getter RpcAttrsGetter[REQUEST]
}

func (r *RpcSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	service := r.getter.GetService(request)
	method := r.getter.GetMethod(request)
	if service == "" || method == "" {
		return "RPC request"
	}
	return service + "/" + method
}
