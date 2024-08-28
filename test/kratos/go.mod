module kratos

go 1.21

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

require (
	github.com/go-kratos/kratos/v2 v2.7.2
	github.com/google/wire v0.6.0
	go.uber.org/automaxprocs v1.5.3
	google.golang.org/genproto/googleapis/api v0.0.0-20240520151616-dc85e6b867a5
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1

)

require (
	github.com/alibaba/opentelemetry-go-auto-instrumentation v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v1.28.0
)

require (
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-fc5f0ca64291 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
