// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fmt7

import (
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func TestOnEnterSprintf3(t *testing.T) {
	ctx := api.NewCallContext()
	arg1 := "format"
	arg2 := "arg"
	ctx.SetParam(1, arg1)
	ctx.SetParam(2, arg1)
	onEnterSprintf3(ctx, arg1, arg2)
}
