package rule

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	api.NewRule("dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown", "Invoke", "*providerGracefulShutdownFilter", "DubboServerOnEnter", "DubboServerOnExit").
		WithFileDeps("dubbo_data_type.go").
		WithFileDeps("dubbo_otel_instrumenter.go").
		WithVersion("[3.0.1,3.1.1)").
		Register()
	api.NewRule("dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown", "Invoke", "*consumerGracefulShutdownFilter", "DubboClientOnEnter", "DubboClientOnExit").
		WithFileDeps("dubbo_data_type.go").
		WithFileDeps("dubbo_otel_instrumenter.go").
		WithVersion("[3.0.1,3.1.1)").
		Register()
}
