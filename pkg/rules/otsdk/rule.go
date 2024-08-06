package otsdk

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "ot_trace_context_linker.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "trace-context/ot_trace_context.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "trace-context/span.go").WithReplace(true).Register()
	api.NewFileRule("go.opentelemetry.io/otel/sdk/trace", "trace-context/tracer.go").WithReplace(true).Register()
	api.NewFileRule("go.opentelemetry.io/otel", "trace-context/trace.go").WithReplace(true).Register()
	// baggage
	api.NewFileRule("go.opentelemetry.io/otel/baggage", "ot_baggage_linker.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/baggage", "ot_baggage_util.go").Register()
	api.NewFileRule("go.opentelemetry.io/otel/baggage", "baggage/context.go").WithReplace(true).Register()
}
