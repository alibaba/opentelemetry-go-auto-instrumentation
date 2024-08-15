//go:build ignore

package rule

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func zapLogWriteOnEnter(call zapcore.CallContext, ce *zapcore.CheckedEntry, fields ...zap.Field) {
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
