// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error17

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onExitBadDep(call api.CallContext, _ string) {
	call.SetReturnVal(0, "gooddep")
	call.SetSkipCall(true)
}
