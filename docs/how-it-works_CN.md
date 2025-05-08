# 工作原理

在正常情况下，`go build` 命令会经历以下几个主要步骤来编译 Golang 应用程序：

1. **源代码解析**：Golang 编译器首先解析源代码文件，并将其转换为抽象语法树（AST）。
2. **类型检查**：解析后，进行类型检查，确保代码符合 Golang 的类型系统。
3. **语义分析**：分析程序的语义，包括变量定义与使用、包导入等。
4. **编译优化**：将语法树转换为中间表示，并进行各种优化，以提高代码执行效率。
5. **代码生成**：生成目标平台的机器码。
6. **链接**：将不同的包和库链接在一起，形成一个单一的可执行文件。

在使用我们的自动注入工具时，会在上述步骤之前新增两个阶段：**预处理（Preprocessing）** 和 **注入（Instrument）**。

![](workflow.png)

- `Preprocess`：分析依赖并选择后续需要使用的注入规则。
- `Instrument`：根据规则生成代码，并将新代码注入到源代码中。

## 预处理阶段（Preprocess）

在此阶段，工具分析用户项目代码中的第三方库依赖，并将其与现有的注入规则进行匹配，找到合适的规则。它还会预先配置这些规则所需要的额外依赖。

注入规则精确定义了哪些代码需要注入到哪个版本的哪个框架或标准库中。不同类型的注入规则有不同的目的。目前可用的注入规则类型包括：

- **InstFuncRule**：在方法的入口和出口插入代码。
- **InstStructRule**：通过添加新字段修改结构体。
- **InstFileRule**：添加新文件参与原有的编译过程。

一旦所有预处理完成，`go build -toolexec otel cmd/app` 会被调用进行编译。`-toolexec` 参数是我们自动注入的核心，用于拦截常规的构建过程并用用户定义的工具替代它，从而让开发者能够更加灵活地定制构建过程。在这里，`otel` 是自动注入工具，这也引出了我们进入注入阶段（Instrument）。

## 注入阶段（Instrument）

在此阶段，根据规则将 trampoline 代码插入到目标函数中。Trampoline 代码本质上是一个复杂的 *If-statement*，它允许在目标函数的入口和出口插入监控代码，从而实现数据采集。此外，还会在 AST 层面进行一些优化，以最小化 trampoline 代码的额外性能开销，并优化代码执行效率。

完成这些步骤后，工具会修改编译参数，然后调用 `go build cmd/app` 进行常规编译，如前所述。

# `net/http` 示例

首先，我们将以下三种类型的函数进行分类：*RawFunc*、*TrampolineFunc*、*HookFunc*。RawFunc 是需要注入的原始函数。TrampolineFunc 是 trampoline 函数。HookFunc 是在原始函数的入口和出口插入的 onEnter/onExit 函数，作为探针代码。RawFunc 通过插入的 trampoline 代码跳转到 TrampolineFunc，然后 TrampolineFunc 构建上下文，准备错误恢复处理，最后跳转到 HookFunc 执行探针代码。

![](tjump.png)

接下来，我们以 `net/http` 为例，演示如何通过编译时的自动注入将监控代码插入到目标函数 `(*Transport).RoundTrip()` 中。框架会在该函数的入口生成 trampoline 代码，这是一条 if 语句（实际上是一行代码，为了演示将其分为多行），它跳转到 TrampolineFunc：

```go
func (t *Transport) RoundTrip(req *Request) (retVal0 *Response, retVal1 error) {
    if callContext37639, _ := OtelOnEnterTrampoline_RoundTrip37639(&t, &req); false {
    } else {
        defer OtelOnExitTrampoline_RoundTrip37639(callContext37639, &retVal0, &retVal1)
    }
    return t.roundTrip(req)
}
```


这里，`OtelOnEnterTrampoline_RoundTrip37639` 是 TrampolineFunc。它准备错误处理和调用上下文，然后跳转到 `ClientOnEnterImpl`：
```go
func OtelOnEnterTrampoline_RoundTrip37639(t **Transport, req **Request) (*CallContext, bool) {
    defer func() {
        if err := recover(); err != nil {
            println("failed to exec onEnter hook", "clientOnEnter")
            if e, ok := err.(error); ok {
                println(e.Error())
            }
            fetchStack, printStack := OtelGetStackImpl, OtelPrintStackImpl
            if fetchStack != nil && printStack != nil {
                printStack(fetchStack())
            }
        }
    }()
    callContext := &CallContext{
        Params:     nil,
        ReturnVals: nil,
        SkipCall:   false,
    }
    callContext.Params = []interface{}{t, req}
    ClientOnEnterImpl(callContext, *t, *req)
    return callContext, callContext.SkipCall
}
```

`ClientOnEnterImpl` 是 HookFunc，它是我们的探针代码，用于执行跟踪、指标报告等操作。`ClientOnEnterImpl` 是一个函数指针，在预处理阶段自动生成的 *otel_setup_inst.go* 文件中进行了预配置，它实际指向 `clientOnEnter`：
```go
// == otel_setup_inst.go
package otel_rules

import http328 "net/http"
...

func init() {
    http328.ClientOnEnterImpl = clientOnEnter
    ...
}
```

`clientOnEnter` 函数执行实际的监控任务：
```go
// == otel_rule_http59729.go
func clientOnEnter(call api.CallContext, t *http.Transport, req *http.Request) {
    ...
    var tracer trace.Tracer
    if span := trace.SpanFromContext(req.Context()); span.SpanContext().IsValid() {
        tracer = span.TracerProvider().Tracer("")
    } else {
        tracer = otel.GetTracerProvider().Tracer("")
    }
    opts := append([]trace.SpanStartOption{}, trace.WithSpanKind(trace.SpanKindClient))
    ctx, span := tracer.Start(req.Context(), req.URL.Path, opts...)
    var attrs []attribute.KeyValue
    attrs = append(attrs, semconv.HTTPMethodKey.String(req.Method))
    attrs = append(attrs, attributes.MakeSpanAttrs(req.URL.Path, req.URL.Host, attributes.Http)...)
    span.SetAttributes(attrs...)
    bag := baggage.FromContext(ctx)
    if mem, err := baggage.NewMemberRaw(constants.BAGGAGE_PARENT_PID, attributes.Pid); err == nil {
        bag, _ = bag.SetMember(mem)
    }
    if mem, err := baggage.NewMemberRaw(constants.BAGGAGE_PARENT_RPC, sdktrace.GetRpc()); err == nil {
        bag, _ = bag.SetMember(mem)
    }
    sdktrace.SetGLocalData(constants.TRACE_ID, span.SpanContext().TraceID().String())
    sdktrace.SetGLocalData(constants.SPAN_ID, span.SpanContext().SpanID().String())
    ctx = baggage.ContextWithBaggage(ctx, bag)
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    req = req.WithContext(ctx)
    *(call.Params[1].(**http.Request)) = req
    return
}
```
通过上述步骤，我们不仅将监控代码插入到 `(*Transport).RoundTrip()` 函数中，还确保了监控数据和上下文的准确性和传播。在编译时自动插桩过程中，这些操作都是自动完成的，从而为开发人员节省了大量时间，并减少了手动探针的错误率。