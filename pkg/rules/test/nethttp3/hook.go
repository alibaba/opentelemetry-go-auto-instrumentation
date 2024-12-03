// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp3

import (
	"context"
	"io"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

// arg type has package prefix
func onEnterNewRequestWithContext(call api.CallContext, ctx context.Context, method, url string, body io.Reader) {
	println("NewRequestWithContext()")
}
