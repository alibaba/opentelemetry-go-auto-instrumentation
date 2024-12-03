// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error14

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onExitOnlyRet(call api.CallContext, _ int, _ string) {
	call.SetReturnVal(0, 2033)
	call.SetReturnVal(1, "hangzhou")
}
