// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fmt5

import (
	_ "fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func OnEnterPrintf2(call api.CallContext, format interface{}, arg ...interface{}) {
	println("hook2")
	for i := 0; i < 10; i++ {
		if i == 5 {
			panic("deliberately")
		}
	}
}
