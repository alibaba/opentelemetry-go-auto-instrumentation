package runtime

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	// Add tls field to struct "g" in runtime2.go
	api.NewStructRule("runtime", "g", "ot_trace_context", "interface{}").
		Register()
	api.NewStructRule("runtime", "g", "ot_baggage_container", "interface{}").
		Register()

	// This solely inspired by skywalking-go
	// https://github.com/apache/skywalking-go/blob/5d7bd5e8e435ec5ab1a61793cd08e6a403893a55/tools/go-agent/instrument/runtime/instrument.go#L75
	api.NewRule("runtime",
		"newproc1", "", "defer func(){ retVal0.ot_trace_context = contextPropagate(callergp.ot_trace_context); retVal0.ot_baggage_container = contextPropagate(callergp.ot_baggage_container); }()", "").
		WithUseRaw(true).
		Register()

	api.NewFileRule("runtime", "runtime_linker.go").
		Register()
}
