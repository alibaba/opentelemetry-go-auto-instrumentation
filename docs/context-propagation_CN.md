## 上下文传递（Context Propagation）

`opentelemetry-go-auto-instrumentation` 中的上下文传递机制受到了 [Apache SkyWalking](https://github.com/apache/skywalking-go) 的启发。

在 OpenTelemetry 中，**上下文（Context）** 是用于在分布式系统中传播追踪相关信息的一种设计。基于上下文的传播，不同的分布式服务（即 Span）可以相互关联，形成完整的调用链（即 Trace）。OpenTelemetry 将追踪相关信息保存在 Golang 的 `context.Context` 中，并要求用户在调用链中正确传递该对象。如果在调用过程中未正确传递 `context.Context`，调用链将会被中断。

为了解决这个问题，当 `opentelemetry-go-auto-instrumentation` 创建一个 Span 时，它会将该 Span 保存到 Golang 的协程结构中（即 GLS）。在创建新协程时，`opentelemetry-go-auto-instrumentation` 还会将当前协程中的相应数据结构一并复制过去。当之后需要创建新的 Span 时，`opentelemetry-go-auto-instrumentation` 会从 GLS 中查询最近创建的 Span 作为其父级，从而有机会保护整个调用链的完整性。

---

Baggage 是 OpenTelemetry 中的一种数据结构，用于在 Trace 中共享键值对。它保存在 `context.Context` 中，并随该对象一同传播。如果 `context.Context` 未在调用链中正确传递，后续服务将无法读取 Baggage。

为了解决这个问题，当 `opentelemetry-go-auto-instrumentation` 将 Baggage 保存到 `context.Context` 时，它也会将 Baggage 保存到 GLS 中。当 `context.Context` 传递不当时，`opentelemetry-go-auto-instrumentation` 将尝试从 GLS 中读取 Baggage，从而使得在这种情况下仍可访问这些数据。