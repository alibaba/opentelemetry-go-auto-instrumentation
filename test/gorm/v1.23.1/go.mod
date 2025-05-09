module gorm

go 1.22

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../../opentelemetry-go-auto-instrumentation/test/verifier



require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.31.0
	go.opentelemetry.io/otel/sdk v1.31.0
	gorm.io/driver/mysql v1.3.1
	gorm.io/gorm v1.23.1
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.31.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.opentelemetry.io/otel v1.31.0
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
