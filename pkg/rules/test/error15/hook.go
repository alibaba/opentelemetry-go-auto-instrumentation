// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error15

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterOnlyArgs(call api.CallContext, _ int, _ string) {
	call.SetParam(0, 2024)
	call.SetParam(1, "shanghai")
}
