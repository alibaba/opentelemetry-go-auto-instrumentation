module gorm

go 1.22

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../../opentelemetry-go-auto-instrumentation/test/verifier

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../opentelemetry-go-auto-instrumentation

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v1.30.0
	go.opentelemetry.io/otel/sdk/metric v1.30.0
	google.golang.org/grpc v1.66.1
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation v0.2.0
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.55.0
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.30.0
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0
	golang.org/x/sys v0.26.0 // indirect
)
