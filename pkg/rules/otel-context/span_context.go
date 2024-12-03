// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel_context

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	spectrace "go.opentelemetry.io/otel/trace"
)

func spanFromContextOnExit(call api.CallContext, span spectrace.Span) {
	if !span.SpanContext().IsValid() {
		call.SetReturnVal(0, sdktrace.SpanFromGLS())
	}
	return
}
