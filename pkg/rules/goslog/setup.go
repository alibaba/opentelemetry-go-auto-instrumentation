// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package golog

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"log/slog"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"go.opentelemetry.io/otel/sdk/trace"
)

var goSlogEnabler = instrumenter.NewDefaultInstrumentEnabler()

func goSlogWriteOnEnter(call api.CallContext, ce *slog.Logger, ctx context.Context, level slog.Level, msg string, args ...any) {
	if !goSlogEnabler.Enable() {
		return
	}
	traceId, spanId := trace.GetTraceAndSpanId()
	if traceId != "" {
		msg = msg + " trace_id=" + traceId
	}
	if spanId != "" {
		msg = msg + " span_id=" + spanId
	}
	call.SetParam(3, msg)
	return
}
