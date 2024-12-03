// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp9

import (
	"net/http"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterClientDo2(call api.CallContext, recv *http.Client, req *http.Request) {
	println("Client.Do2()")
}
