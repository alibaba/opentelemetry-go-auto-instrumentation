module build

go 1.22

toolchain go1.23.0

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

require (
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
)
