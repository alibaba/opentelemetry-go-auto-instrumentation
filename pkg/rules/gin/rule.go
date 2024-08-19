package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {
	api.NewRule("github.com/gin-gonic/gin", "Next", "*Context", "nextOnEnter", "nextOnExit").
		WithVersion("[1.7.0,1.10.1)").
		WithFileDeps("gin_data_type.go").
		WithFileDeps("gin_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/gin-gonic/gin", "HTML", "*Context", "htmlOnEnter", "htmlOnExit").
		WithVersion("[1.7.0,1.10.1)").
		WithFileDeps("gin_data_type.go").
		WithFileDeps("gin_otel_instrumenter.go").
		Register()
}
