module github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gorm

go 1.23.0

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg => ../../../pkg

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.9.3
	go.opentelemetry.io/otel/sdk v1.36.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.22.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
