package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {

	//client
	api.NewRule("net/http", "RoundTrip", "*Transport", "clientOnEnter", "clientOnExit").
		WithFileDeps("net_http_data_type.go").
		WithFileDeps("net_http_otel_instrumenter.go").
		Register()

	//server
	api.NewRule("net/http", "ServeHTTP", "serverHandler", "serverOnEnter", "serverOnExit").
		WithFileDeps("net_http_data_type.go").
		WithFileDeps("net_http_otel_instrumenter.go").
		Register()

}
