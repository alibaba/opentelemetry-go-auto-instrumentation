//go:build ignore

package rule

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

func logNewOnExit(call logrus.CallContext) {
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
