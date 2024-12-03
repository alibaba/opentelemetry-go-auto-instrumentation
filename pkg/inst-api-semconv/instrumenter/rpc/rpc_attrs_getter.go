// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

type RpcAttrsGetter[REQUEST any] interface {
	GetSystem(request REQUEST) string
	GetService(request REQUEST) string
	GetMethod(request REQUEST) string
}
