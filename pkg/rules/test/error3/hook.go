// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error2

import (
	erralias "errors"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterErrorsNewAlias(call api.CallContext, text string) {
	// Check if alias name confuses compilation
	_ = erralias.ErrUnsupported
}
