// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp7

import (
	"net/http"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

// use field added by struct rule
func onExitNewRequest(call api.CallContext, req *http.Request, _ interface{}) {
	println(req.Should)
}
