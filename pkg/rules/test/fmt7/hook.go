// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fmt7

import (
	_ "fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterSprintf3(call api.CallContext, format string, arg ...any) {
	println("a3")
}

func onExitSprintf3(call api.CallContext, s string) {
	print("b3")
}
