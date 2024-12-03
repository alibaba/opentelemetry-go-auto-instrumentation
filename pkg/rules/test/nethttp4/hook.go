// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp4

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

// many args have interface type
func onEnterNewRequest1(call api.CallContext, a, b interface{}, c interface{}) {
	println("NewRequest1()")
}
