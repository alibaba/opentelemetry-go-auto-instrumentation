// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp5

import (
	"net/http"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

// only recv arg
func onEnterMaxBytesError(call api.CallContext, recv *http.MaxBytesError) {
	println("MaxBytesError()")
	recv.Limit = 4008208820
}

func onExitMaxBytesError(call api.CallContext, ret string) {
	call.SetReturnVal(0, "Prince of Qin Smashing the Battle line")
}
