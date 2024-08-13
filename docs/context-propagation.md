## Context Propagation

Context Propagation in `opentelemetry-go-auto-instrumentation` is inspired
by [Apache-Skywalking](https://github.com/apache/skywalking-go).
Context in OpenTelemetry is a design for propagating trace-related information in distributed systems. Based on the
propagation of Context, distributed services (AKA Spans) can be linked together to form a complete call chain (AKA
Trace). OpenTelemetry saves trace-related information in Golang's context.Context and requires users to correctly pass
context.Context. If context.Context is not passed correctly in the call chain, the call chain will be broken. To solve
this problem, when `opentelemetry-go-auto-instrumentation` create a span, `opentelemetry-go-auto-instrumentation` save
it to Golang's coroutine structure (i.e. GLS), and when `opentelemetry-go-auto-instrumentation` create a new
coroutine, `opentelemetry-go-auto-instrumentation` also copy the corresponding data structure from the current
coroutine. When `opentelemetry-go-auto-instrumentation` need to create a new span
later, `opentelemetry-go-auto-instrumentation` will query the most recently created span from GLS as the parent, so that
`opentelemetry-go-auto-instrumentation` have the opportunity to protect
the integrity of the call chain.

Baggage is a data structure in OpenTelemetry used to share key-value pairs in Trace. Baggage is stored in
context.Context and is propagated along with the context.Context. If context.Context is not correctly propagated in the
call chain, subsequent services will not be able to read Baggage. To solve this problem,
when `opentelemetry-go-auto-instrumentation` save the baggage to
context.Context, `opentelemetry-go-auto-instrumentation` also save it to GLS. When context.Context is not passed
correctly, `opentelemetry-go-auto-instrumentation` try to read the baggage from
GLS, which allows the baggage to be read in this case.