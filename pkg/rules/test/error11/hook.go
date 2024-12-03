// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error11

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterTestSkip2(call api.CallContext) {
	call.SetSkipCall(true)
}

func onExitTestSkip2(call api.CallContext, _ int) {
	call.SetReturnVal(0, 0x512)
}
