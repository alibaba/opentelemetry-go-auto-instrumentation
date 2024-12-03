// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fasthttp

import (
	"net/url"

	"github.com/valyala/fasthttp"
)

type fastHttpRequest struct {
	method string
	url    *url.URL
	isTls  bool
	port   int
	header *fasthttp.RequestHeader
}

type fastHttpResponse struct {
	statusCode int
	header     *fasthttp.ResponseHeader
}
