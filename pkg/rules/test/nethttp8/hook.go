// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nethttp8

import (
	"context"
	"io"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterNewRequestWithContext2(call api.CallContext, ctx context.Context, method, url string, body io.Reader) {
	println("NewRequestWithContext2()")
}
