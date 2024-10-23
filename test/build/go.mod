module build

go 1.22

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240822170219-fc7c04adadcd

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

require (
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
)
