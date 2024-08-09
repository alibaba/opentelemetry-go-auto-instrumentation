package runtime

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	// Add tls field to struct "g" in runtime2.go
	api.NewStructRule("runtime", "g", "ot_trace_context", "interface{}").
		Register()
	api.NewStructRule("runtime", "g", "ot_baggage_container", "interface{}").
		Register()
	// Defer call
	api.NewRule("runtime",
		"newproc1", "", "defer func(){ retVal0.ot_trace_context = contextPropagate(callergp.ot_trace_context); retVal0.ot_baggage_container = contextPropagate(callergp.ot_baggage_container); }()", "").
		WithUseRaw(true).
		Register()

	api.NewFileRule("runtime", "runtime_linker.go").
		Register()
}
