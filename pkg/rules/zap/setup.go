//go:build ignore

package rule

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func zapLogWriteOnEnter(call zapcore.CallContext, ce *zapcore.CheckedEntry, fields ...zap.Field) {
	var fieldsTemp []zap.Field
	traceId, spanId := trace.GetGLocalDataDouble("trace_id", "span_id")
	if traceId != nil {
		fieldsTemp = append(fieldsTemp, zap.String("trace_id", traceId.(string)))
	}
	if spanId != nil {
		fieldsTemp = append(fieldsTemp, zap.String("span_id", spanId.(string)))
	}
	if fields == nil {
		fields = fieldsTemp
	}
	call.SetParam(1, fields)
	return
}
