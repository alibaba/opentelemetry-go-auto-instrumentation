// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package golog

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"log"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"go.opentelemetry.io/otel/sdk/trace"
)

var glogEnabler = instrumenter.NewDefaultInstrumentEnabler()

func goLogWriteOnEnter(call api.CallContext, ce *log.Logger, pc uintptr, calldepth int, appendOutput func([]byte) []byte) {
	if !glogEnabler.Enable() {
		return
	}
	traceId, spanId := trace.GetTraceAndSpanId()
	newAppendOutput := func(bytes []byte) []byte {
		sb := strings.Builder{}
		if traceId != "" {
			sb.WriteString(" trace_id=")
			sb.WriteString(traceId)
		}
		if spanId != "" {
			sb.WriteString(" span_id=")
			sb.WriteString(spanId)
		}
		bytes = append(bytes, []byte(sb.String())...)
		bytes = appendOutput(bytes)
		sb.Reset()
		return bytes
	}
	call.SetParam(3, newAppendOutput)
	return
}
