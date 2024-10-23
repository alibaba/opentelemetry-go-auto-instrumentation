module grpc

go 1.22.7

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/verifier => ../../../opentelemetry-go-auto-instrumentation/verifier

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation/verifier v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v1.30.0
	google.golang.org/grpc v1.66.2
	google.golang.org/grpc/examples v0.0.0-20240925225118-218811eb43b1
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
)
