package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {
	api.NewRule("github.com/labstack/echo/v4", "New", "", "", "afterNewEcho").
		WithVersion("[4.0.0,4.12.1)").
		WithFileDeps("echo_data_type.go", "echo_otel_instrumenter.go").
		Register()
}
