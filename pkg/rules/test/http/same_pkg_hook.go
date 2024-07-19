//go:build ignore

package http

import (
	"context"
	"io"
	"net/http"
)

func onEnterClientDo2(call *http.CallContext, recv *http.Client, req *http.Request) {
	println("Client.Do2()")
}

func onEnterNewRequestWithContext2(call *http.CallContext, ctx context.Context, method, url string, body io.Reader) {
	println("NewRequestWithContext2()")
}
