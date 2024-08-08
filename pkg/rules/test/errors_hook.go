//go:build ignore

package test

import (
	"errors"
	"fmt"
)

func onEnterUnwrap(call errors.CallContext, err error) {
	newErr := fmt.Errorf("wrapped: %w", err)
	call.SetParam(0, newErr)
}

func onExitUnwrap(call errors.CallContext, err error) {
	e := call.GetParam(0).(interface {
		Unwrap() error
	})
	old := e.Unwrap()
	fmt.Printf("old:%v\n", old)
}

func onEnterTestSkip(call errors.CallContext) {
	call.SetSkipCall(true)
}

func onExitTestSkipOnly(call errors.CallContext, _ *int) {}

func onEnterTestSkipOnly(call errors.CallContext) {}

func onEnterP11(call errors.CallContext) {}
func onEnterP12(call errors.CallContext) {}

func onExitP21(call errors.CallContext) {}
func onExitP22(call errors.CallContext) {}

func onEnterP31(call errors.CallContext, arg1 int, arg2 bool, arg3 float64) {}
func onExitP31(call errors.CallContext, arg1 int, arg2 bool, arg3 float64)  {}

func onEnterTestSkip2(call errors.CallContext) {
	call.SetSkipCall(true)
}

func onExitTestSkip2(call errors.CallContext, _ int) {
	call.SetReturnVal(0, 0x512)
}

func onEnterTestGetSet(call errors.CallContext, arg1 int, arg2, arg3 bool, arg4 float64, arg5 string, arg6 interface{}, arg7, arg8 map[int]bool, arg9 chan int, arg10 []int) {
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

func onExitTestGetSet(call errors.CallContext, arg1 int, arg2 bool, arg3 bool, arg4 float64, arg5 string, arg6 interface{}, arg7 map[int]bool, arg8 map[int]bool, arg9 chan int, arg10 []int) {
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

func onEnterTestGetSetRecv(call errors.CallContext, arg1 *errors.Recv, arg2 int, arg3 float64) {
	recv := call.GetParam(0).(*errors.Recv)
	recv.X = 4008208820
	// call.SetParam(0, recv)
	call.SetParam(1, 118888)
	call.SetParam(2, 0.001)
}

func onExitTestGetSetRecv(call errors.CallContext, arg1 int, arg2 float64) {
	call.SetReturnVal(0, arg1)
	call.SetReturnVal(1, arg2)
}

func onExitOnlyRet(call errors.CallContext, _ int, _ string) {
	call.SetReturnVal(0, 2033)
	call.SetReturnVal(1, "hangzhou")
}

func onEnterOnlyArgs(call errors.CallContext, _ int, _ string) {
	call.SetParam(0, 2024)
	call.SetParam(1, "shanghai")
}

func onEnterNilArg(call errors.CallContext, _ *int) {
	// GetParam(0) is nil
	arg0 := call.GetParam(0)
	println("getparam0", arg0)
	call.SetParam(0, nil)
}
