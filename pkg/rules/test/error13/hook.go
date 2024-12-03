// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error13

import (
	"reflect"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterTestGetSetRecv(call api.CallContext, arg1 interface{}, arg2 int, arg3 float64) {
	recv := call.GetParam(0)
	v := reflect.ValueOf(recv)
	v = v.Elem()
	field := v.FieldByName("X")
	field.SetInt(4008208820)

	call.SetParam(1, 118888)
	call.SetParam(2, 0.001)
}

func onExitTestGetSetRecv(call api.CallContext, arg1 int, arg2 float64) {
	call.SetReturnVal(0, arg1)
	call.SetReturnVal(1, arg2)
}
