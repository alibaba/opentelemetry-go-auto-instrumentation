// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fiberv2

import (
	"net/url"

	"github.com/valyala/fasthttp"
)

type fiberv2Request struct {
	method string
	url    *url.URL
	isTls  bool
	port   int
	header *fasthttp.RequestHeader
}

type fiberv2Response struct {
	statusCode int
	header     *fasthttp.ResponseHeader
}
