// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

type RpcSpanNameExtractor[REQUEST any] struct {
	Getter RpcAttrsGetter[REQUEST]
}

func (r *RpcSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	service := r.Getter.GetService(request)
	method := r.Getter.GetMethod(request)
	if service == "" || method == "" {
		return "RPC request"
	}
	return service + "/" + method
}
