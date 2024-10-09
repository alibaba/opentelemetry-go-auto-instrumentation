## How to debug instrumented program

`opentelemetry-go-auto-instrumentation` provides some convenient ways for users to debug the instrumented program.

## 1. Perform instrumentation with -debug options

```bash
$ ./otelbuild -debug
```

When using the `-debug` compilation option, the tool will compile an unoptimized binary 
while retaining all generated temporary files, such as `otel_rules`. You can review 
them to understand what kind of code the tool is injecting.

## 2. Check `.otel-build` directory

Even without using the `-debug` option, the tool will retain the necessary modified file copies in `.otel-build`, and its structure is as follows:

```shell
.otel-build
├── instrument
│   ├── debug_file_otel_inst_file_context.go
│   ├── debug_file_otel_inst_file_ot_baggage_linker.go
│   ├── debug_file_otel_inst_file_ot_baggage_util.go
│   ├── debug_file_otel_inst_file_ot_trace_context.go
│   ├── debug_file_otel_inst_file_ot_trace_context_linker.go
│   ├── debug_file_otel_inst_file_runtime_linker.go
│   ├── debug_file_otel_inst_file_span.go
│   ├── debug_file_otel_inst_file_trace.go
│   ├── debug_file_otel_inst_file_tracer.go
│   ├── debug_fn_clientconn.go
│   ├── debug_fn_entry.go
│   ├── debug_fn_proc.go
│   ├── debug_fn_roundtrip.go
│   ├── debug_fn_server.go
│   ├── debug_grpc_otel_trampoline.go
│   ├── debug_http_otel_trampoline.go
│   ├── debug_runtime_otel_trampoline.go
│   ├── debug_struct_runtime2.go
│   └── debug_zapcore_otel_trampoline.go
└── preprocess
    ├── backups
    │   ...
    ├── debug_otel_rule_grpc01449.go
    ├── debug_otel_rule_grpc03361.go
    ├── debug_otel_rule_grpc22680.go
    ├── debug_otel_rule_grpc27355.go
    ├── debug_otel_rule_grpc77619.go
    ├── debug_otel_rule_grpc78906.go
    ├── debug_otel_rule_grpc90364.go
    ├── debug_otel_rule_grpc94639.go
    ├── debug_otel_rule_http15118.go
    ├── debug_otel_rule_http39595.go
    ├── debug_otel_rule_http44021.go
    ├── debug_otel_rule_http59627.go
    ├── debug_otel_rule_zapcore35871.go
    ├── debug_otel_setup_inst.go
    ├── debug_otel_setup_sdk.go
    ├── debug_user_main.go
    ├── dry_run.log
    ├── embededfs
    │   ...
    └── used_rules.json
```

The terms "preprocess" and "instrument" represent files generated during two different stages. Please refer to [this document](how-it-works.md) for information about the two stages. For example, `debug_fn_clientconn.go` indicates the `clientconn.go` file after code injection, while `debug_otel_rule_grpc01449.go` refers to the specific hook code. `used_rules.json` contains the matched rules, and nearly all important files relevant to debugging will be retained in this directory.

## 3. Use delve to debug binary

No optimization will be taken with the `-debug` option during the hybrid compilation. Users can
use [delve](https://github.com/go-delve/delve) to debug the binary file easily.