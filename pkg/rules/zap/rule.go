package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {
	api.NewRule("go.uber.org/zap/zapcore", "Write", "*CheckedEntry", "zapLogWriteOnEnter", "").
		WithVersion("[1.20.0,1.27.1)").
		Register()
}
