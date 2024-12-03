// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package echo

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	echo "github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/sdk/trace"
)

var echoEnabler = instrumenter.NewDefaultInstrumentEnabler()

func otelTraceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if err = next(c); err != nil {
				c.Error(err)
			}
			lcs := trace.LocalRootSpanFromGLS()
			if lcs != nil && c.Path() != "" && c.Request() != nil && c.Request().URL != nil && (c.Request().URL.Path != c.Path()) {
				lcs.SetName(c.Path())
			}
			return
		}
	}
}

func afterNewEcho(call api.CallContext, e *echo.Echo) {
	if !echoEnabler.Enable() {
		return
	}
	if e == nil {
		return
	}

	e.Use(otelTraceMiddleware())
}
