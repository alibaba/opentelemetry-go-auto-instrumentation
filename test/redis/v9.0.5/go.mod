module redis/v9.0.5

go 1.21

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../opentelemetry-go-auto-instrumentation

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.0.5
	go.opentelemetry.io/otel/sdk v1.28.0
)