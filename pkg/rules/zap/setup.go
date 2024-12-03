// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package zap

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapEnabler = instrumenter.NewDefaultInstrumentEnabler()

func zapLogWriteOnEnter(call api.CallContext, ce *zapcore.CheckedEntry, fields ...zap.Field) {
	if !zapEnabler.Enable() {
		return
	}
	var fieldsTemp []zap.Field
	traceId, spanId := trace.GetTraceAndSpanId()
	if traceId != "" {
		fieldsTemp = append(fieldsTemp, zap.String("trace_id", traceId))
	}
	if spanId != "" {
		fieldsTemp = append(fieldsTemp, zap.String("span_id", spanId))
	}
	if fields == nil {
		fields = fieldsTemp
	}
	call.SetParam(1, fields)
	return
}
