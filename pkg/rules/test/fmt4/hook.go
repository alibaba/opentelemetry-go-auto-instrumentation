// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fmt4

import (
	_ "fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

type any = interface{}

func onEnterSprintf1(call api.CallContext, format string, arg ...any) {
	print("a1")
}

func onExitSprintf1(call api.CallContext, s string) {
	print("b1")
}
