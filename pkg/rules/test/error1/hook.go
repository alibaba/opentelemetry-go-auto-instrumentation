// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package error1

import (
	_ "errors"
	"fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterUnwrap(call api.CallContext, err error) {
	newErr := fmt.Errorf("wrapped: %w", err)
	call.SetParam(0, newErr)
}

func onExitUnwrap(call api.CallContext, err error) {
	e := call.GetParam(0).(interface {
		Unwrap() error
	})
	old := e.Unwrap()
	fmt.Printf("old:%v\n", old)
}
