// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fmt1

import (
	_ "fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func OnExitPrintf1(call api.CallContext, n int, err error) {
	println("Exiting hook1....")
	call.SetReturnVal(0, 1024)
	v := call.GetData().(int)
	println(v)
}

func OnEnterPrintf1(call api.CallContext, format string, arg ...any) {
	println("Entering hook1....")
	call.SetData(555)
	call.SetParam(0, "olleH%s\n")
	p1 := call.GetParam(1).([]any)
	p1[0] = "goodcatch"
}
