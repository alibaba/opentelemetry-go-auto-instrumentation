module gorm/v1.22.0

go 1.22

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier => ../../../../opentelemetry-go-auto-instrumentation

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v1.30.0
	gorm.io/driver/mysql v1.1.3
	gorm.io/gorm v1.22.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.3 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
