module redigo

go 1.22

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../../opentelemetry-go-auto-instrumentation/test/verifier

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../opentelemetry-go-auto-instrumentation

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation v0.0.0-00010101000000-000000000000
    github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier v0.0.0-00010101000000-000000000000
	github.com/gomodule/redigo v1.9.0
)