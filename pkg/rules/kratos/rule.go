package rule

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

// this plugin can NOT work with rules/http/arms
// DISABLE EITHER OF THEM
func init() {

	api.NewRule("github.com/go-kratos/kratos/v2/transport/http", "NewServer", "", "KratosNewHTTPServiceOnEnter", "").
		WithVersion("[2.5.2,2.7.4)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/go-kratos/kratos/v2/transport/grpc", "NewServer", "", "KratosNewGRPCServiceOnEnter", "").
		WithVersion("[2.5.2,2.7.4)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/go-kratos/kratos/v2/transport/http", "NewClient", "", "KratosNewHTTPClientOnEnter", "").
		WithVersion("[2.5.2,2.7.4)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/go-kratos/kratos/v2/transport/http", "WithMiddleware", "", "KratosWithMiddlewareOnEnter", "").
		WithVersion("[2.5.2,2.7.4)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/go-kratos/kratos/v2/transport/grpc", "DialInsecure", "", "KratosDialInsecureOnEnter", "KratosDialInsecureOnExit").
		WithVersion("[2.5.2,2.7.4)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/go-kratos/kratos/v2/transport/grpc", "WithMiddleware", "", "KratosGRPCWithMiddlewareOnEnter", "").
		WithVersion("[2.5.2,2.7.4)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

}
