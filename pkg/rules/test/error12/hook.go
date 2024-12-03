// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error12

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterTestGetSet(call api.CallContext, arg1 int, arg2, arg3 bool, arg4 float64, arg5 string, arg6 interface{}, arg7, arg8 map[int]bool, arg9 chan int, arg10 []int) {
	call.SetParam(0, 7632)
	call.SetParam(1, arg2)
	call.SetParam(2, arg3)
	call.SetParam(3, arg4)
	call.SetParam(4, arg5)
	call.SetParam(5, arg6)
	call.SetParam(6, arg7)
	call.SetParam(7, arg8)
	call.SetParam(8, arg9)
	call.SetParam(9, arg10)
}

func onExitTestGetSet(call api.CallContext, arg1 int, arg2 bool, arg3 bool, arg4 float64, arg5 string, arg6 interface{}, arg7 map[int]bool, arg8 map[int]bool, arg9 chan int, arg10 []int) {
	call.SetReturnVal(0, arg1)
	call.SetReturnVal(1, arg2)
	call.SetReturnVal(2, arg3)
	call.SetReturnVal(3, arg4)
	call.SetReturnVal(4, arg5)
	call.SetReturnVal(5, arg6)
	call.SetReturnVal(6, arg7)
	call.SetReturnVal(7, arg8)
	call.SetReturnVal(8, arg9)
	call.SetReturnVal(9, arg10)
}
