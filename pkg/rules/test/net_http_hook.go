//go:build ignore

package test

import (
	"context"
	"io"
	"net/http"
	"time"
)

func onEnterClientDo(call *http.CallContext, recv *http.Client, req *http.Request) {
	if time.Now().UnixMilli()%2 == 1 {
		panic("test panic")
	} else {
		println("Before Client.Do()")
	}
}

func onExitClientDo(call *http.CallContext, resp *http.Response, err error) {
	if time.Now().UnixMilli()%2 == 1 {
		println("After Client.Do()")
	} else {
		panic("deliberately")
	}
}

// arg type has package prefix
func onEnterNewRequestWithContext(call *http.CallContext, ctx context.Context, method, url string, body io.Reader) {
	println("NewRequestWithContext()")
}

// many args have one type
func onEnterNewRequest(call *http.CallContext, method, url string, body io.Reader) {
	println("NewRequest()")
}

// many args have interface type
func onEnterNewRequest1(call *http.CallContext, a, b interface{}, c interface{}) {
	println("NewRequest1()")
}

// only recv arg
func onEnterMaxBytesError(call *http.CallContext, recv *http.MaxBytesError) {
	println("MaxBytesError()")
	recv.Limit = 4008208820
}

func onExitMaxBytesError(call *http.CallContext, ret string) {
	*(call.ReturnVals[0].(*string)) = "Prince of Qin Smashing the Battle line"
}

// use field added by struct rule
func onExitNewRequest(call *http.CallContext, req *http.Request, _ interface{}) {
	println(req.Should)
}
