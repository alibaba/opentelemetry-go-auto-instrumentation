package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {
	api.NewRule("github.com/gorilla/mux", "ServeHTTP", "*Router", "muxServerOnEnter", "muxServerOnExit").
		WithVersion("[1.3.0,1.8.2)").
		WithFileDeps("mux_data_type.go", "mux_otel_instrumenter.go").
		Register()
}
