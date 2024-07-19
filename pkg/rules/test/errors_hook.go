//go:build ignore

package test

import (
	"errors"
	"fmt"
)

func onEnterUnwrap(call *errors.CallContext, err error) {
	newErr := fmt.Errorf("wrapped: %w", err)
	*(call.Params[0].(*error)) = newErr
}

func onExitUnwrap(call *errors.CallContext, err error) {
	e := (*(call.Params[0].(*error))).(interface {
		Unwrap() error
	})
	old := e.Unwrap()
	fmt.Printf("old:%v\n", old)
}

func onEnterTestSkip(call *errors.CallContext) {
	call.SetSkipCall(true)
}

func onExitTestSkipOnly(call *errors.CallContext, _ *int) {}

func onEnterTestSkipOnly(call *errors.CallContext) {}

func onEnterP11(call *errors.CallContext) {}
func onEnterP12(call *errors.CallContext) {}

func onExitP21(call *errors.CallContext) {}
func onExitP22(call *errors.CallContext) {}

func onEnterP31(call *errors.CallContext, arg1 int, arg2 bool, arg3 float64) {}
func onExitP31(call *errors.CallContext, arg1 int, arg2 bool, arg3 float64)  {}

func onEnterTestSkip2(call *errors.CallContext) {
	call.SetSkipCall(true)
}

func onExitTestSkip2(call *errors.CallContext, _ int) {
	*(call.ReturnVals[0].(*int)) = 0x512
}
