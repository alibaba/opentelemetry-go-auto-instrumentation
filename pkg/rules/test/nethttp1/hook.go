// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp1

import (
	"net/http"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterClientDo(call api.CallContext, recv *http.Client, req *http.Request) {
	println("Before Client.Do()")
}

func onExitClientDo(call api.CallContext, resp *http.Response, err error) {
	panic("deliberately")
}
