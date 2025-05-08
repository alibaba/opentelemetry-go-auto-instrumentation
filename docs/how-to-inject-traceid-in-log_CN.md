# 如何在日志中注入 traceID

## 自动注入

如果我们使用的是 `opentelemetry-go-auto-instrumentation` 支持的日志框架（参见 [这里](./supported-libraries.md)），则 TraceId 和 SpanId 会被自动注入到日志中。
```go
package main

import (
	"go.uber.org/zap"
	"net/http"
)

func main() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		logger := zap.NewExample()
		logger.Debug("this is debug message")
		logger.Info("this is info message")
		logger.Warn("this is warn message")
		logger.Error("this is error message")
	})

	http.ListenAndServe(":9999", nil)
}

```

例如，如果我们使用 `opentelemetry-go-auto-instrumentation` 编译以下 Go 文件，运行生成的二进制文件并执行 `curl localhost:9999/log`，我们将看到如下输出：

```shell
{"level":"debug","msg":"this is debug message","trace_id":"d62a8fea286cc66de9c68ca17d4faa88","span_id":"7cb6d692769ffd32"}
{"level":"info","msg":"this is info message","trace_id":"d62a8fea286cc66de9c68ca17d4faa88","span_id":"7cb6d692769ffd32"}
{"level":"warn","msg":"this is warn message","trace_id":"d62a8fea286cc66de9c68ca17d4faa88","span_id":"7cb6d692769ffd32"}
{"level":"error","msg":"this is error message","trace_id":"d62a8fea286cc66de9c68ca17d4faa88","span_id":"7cb6d692769ffd32"}
{"level":"debug","msg":"this is debug message","trace_id":"e56a6f1e7ed7af48cce8f64d045ed158","span_id":"def0b8cf10fe8844"}
{"level":"info","msg":"this is info message","trace_id":"e56a6f1e7ed7af48cce8f64d045ed158","span_id":"def0b8cf10fe8844"}
{"level":"warn","msg":"this is warn message","trace_id":"e56a6f1e7ed7af48cce8f64d045ed158","span_id":"def0b8cf10fe8844"}
{"level":"error","msg":"this is error message","trace_id":"e56a6f1e7ed7af48cce8f64d045ed158","span_id":"def0b8cf10fe8844"}
```

TraceId 和 SpanId 会被自动注入到日志中。

## 手动注入

如果使用的日志框架不在 `opentelemetry-go-auto-instrumentation` 支持范围内，我们可以手动将 TraceId 和 SpanId 注入到日志中：
```go
package main

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	http.HandleFunc("/logwithtrace", func(w http.ResponseWriter, r *http.Request) {
		logger := zap.NewExample()
		traceId, spanId := trace.GetTraceAndSpanId()
		logger.Info("this is info message with fileds",
			zap.String("traceId", traceId),
			zap.String("spanId", spanId),
		)
	})
	http.ListenAndServe(":9999", nil)
}
```

例如，如果我们使用 `opentelemetry-go-auto-instrumentation` 构建以下 Go 文件，运行生成的可执行文件后执行 `curl localhost:9999/logwithtrace`，将会看到如下输出：
```shell
{"level":"info","msg":"this is info message with fileds","traceId":"92d63797010a2040484222a74c5ce304","spanId":"5a2c84c807a6e12c"}
```
上述代码位于 [example/log](../example/log) 目录中。