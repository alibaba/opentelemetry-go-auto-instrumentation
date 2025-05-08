## 如何调试已仪表化的程序

`opentelemetry-go-auto-instrumentation` 提供了一些方便的方式供用户调试已仪表化的程序。

## 1. 使用 -debug 选项进行仪表化
```console
$ ./otel set -debug
```

使用 `-debug` 编译选项时，工具会生成一个未优化的二进制文件，并保留所有生成的临时文件，例如 `otel_rules`。你可以查看这些文件，以了解工具注入了哪些代码。

## 2. 查看 `.otel-build` 目录

即使未使用 `-debug` 选项，工具也会在 `.otel-build` 目录中保留必要的已修改文件副本，其结构如下：
```shell
.otel-build
├── instrument
│   ├── baggage
│   │   ├── otel_inst_file_context.go
│   │   ├── otel_inst_file_ot_baggage_linker.go
│   │   └── otel_inst_file_ot_baggage_util.go
│   ├── grpc
│   │   ├── clientconn.go
│   │   ├── otel_trampoline.go
│   │   └── server.go
│   ├── http
│   │   ├── otel_trampoline.go
│   │   ├── roundtrip.go
│   │   └── server.go
│   ├── log
│   │   ├── log.go
│   │   └── otel_trampoline.go
│   ├── otel
│   │   └── otel_inst_file_trace.go
│   ├── runtime
│   │   ├── otel_inst_file_runtime_linker.go
│   │   ├── otel_trampoline.go
│   │   ├── proc.go
│   │   └── runtime2.go
│   ├── slog
│   │   ├── logger.go
│   │   └── otel_trampoline.go
│   └── trace
│       ├── otel_inst_file_ot_trace_context.go
│       ├── otel_inst_file_ot_trace_context_linker.go
│       ├── otel_inst_file_span.go
│       └── otel_inst_file_tracer.go
└── preprocess
    ├── backups
    │   ├── app2.go.bk
    │   ├── go.mod.bk
    │   └── go.sum.bk
    ├── dry_run.log
    ├── otel_rules
    │   ├── grpc72047
    │   │   ├── ...
    │   ├── http02075
    │   │   ├── client_setup.go
    │   │   ├── ...
    │   ├── log09344
    │   │   └── setup.go
    │   ├── otel_setup_inst.go
    │   ├── otel_setup_sdk.go
    │   └── slog54146
    │       └── setup.go
    ├── otel_user
    │   ├── app2.go
    │   ├── go.mod
    │   └── go.sum
    ├── rule_bundle.json
    └── rule_cache
        └── ...
```

术语 “preprocess” 和 “instrument” 代表了在两个不同阶段生成的文件。关于这两个阶段的详细信息，请参考 [这篇文档](how-it-works_CN.md)。例如，`instrument/grpc/clientconn.go` 表示注入代码后的 `clientconn.go` 文件。`rule_bundle.json` 包含了匹配到的规则，几乎所有与调试相关的重要文件都会保留在此目录中。

## 3. 使用 delve 调试二进制文件

在使用 `-debug` 选项进行混合编译时，工具不会进行任何优化。用户可以使用 [delve](https://github.com/go-delve/delve) 来方便地调试该二进制文件。