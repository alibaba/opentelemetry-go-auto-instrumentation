module otel

go 1.23.0

toolchain go1.24.2

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../opentelemetry-go-auto-instrumentation/test/verifier

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

require go.opentelemetry.io/otel/trace v1.35.0

require go.opentelemetry.io/otel v1.35.0 // indirect
