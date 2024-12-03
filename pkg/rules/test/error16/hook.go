// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error16

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterNilArg(call api.CallContext, _ *int) {
	// GetParam(0) is nil
	arg0 := call.GetParam(0)
	println("getparam0", arg0)
	call.SetParam(0, nil)
}
