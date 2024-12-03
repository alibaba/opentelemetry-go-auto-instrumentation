// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp2

import (
	"io"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

// many args have one type
func onEnterNewRequest(call api.CallContext, method, url string, body io.Reader) {
	println("NewRequest()")
}
