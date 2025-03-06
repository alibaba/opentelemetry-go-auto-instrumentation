module langchain/v0.1.13

go 1.22.0

toolchain go1.23.3

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../opentelemetry-go-auto-instrumentation

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../../opentelemetry-go-auto-instrumentation/test/verifier

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier v0.0.0-20250227024304-f25d514205ee
	github.com/tmc/langchaingo v0.1.13
	go.opentelemetry.io/otel/sdk v1.34.0
)

require (
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pkoukk/tiktoken-go v0.1.6 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
