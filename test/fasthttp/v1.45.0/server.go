// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/valyala/fasthttp"

func hello(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("hello world"))
}
