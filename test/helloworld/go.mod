module helloworld

go 1.22.0

toolchain go1.22.7

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../opentelemetry-go-auto-instrumentation/test/verifier

require (
	go.opentelemetry.io/otel v1.33.0
	golang.org/x/time v0.5.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
)
