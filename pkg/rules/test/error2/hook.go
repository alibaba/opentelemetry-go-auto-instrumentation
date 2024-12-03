// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error2

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterTestSkip(call api.CallContext) {
	call.SetSkipCall(true)
}
