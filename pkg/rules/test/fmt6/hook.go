// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fmt6

import (
	_ "fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterSprintf2(call api.CallContext, format string, arg ...any) {
	print("a2")
	_ = call.IsSkipCall()
}

func onExitSprintf2(call api.CallContext, s string) {
	println("b2")
}
