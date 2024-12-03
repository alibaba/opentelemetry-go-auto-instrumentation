// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logrus

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

var logrusEnabler = instrumenter.NewDefaultInstrumentEnabler()

func logNewOnExit(call api.CallContext) {
	if !logrusEnabler.Enable() {
		return
	}
	std := logrus.StandardLogger()
	std.AddHook(&logHook{})
	return
}

type logHook struct{}

func (hook *logHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *logHook) Fire(entry *logrus.Entry) error {
	// 修改日志内容
	traceId, spanId := trace.GetTraceAndSpanId()
	if traceId != "" {
		entry.Data["trace_id"] = traceId
	}
	if spanId != "" {
		entry.Data["span_id"] = spanId
	}
	return nil
}
