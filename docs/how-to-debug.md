## How to debug instrumented program

`loongsuite-go-agent` provides some convenient ways for users to debug the instrumented program.

## 1. Perform instrumentation with -debug options

```console
$ ./otel set -debug
```

When using the `-debug` compilation option, the tool will compile an unoptimized binary 
while retaining all generated temporary files, such as `otel_rules`. You can review 
them to understand what kind of code the tool is injecting.

## 2. Check `.otel-build` directory

Even without using the `-debug` option, the tool will retain the necessary modified file copies in `.otel-build`, and its structure is as follows:

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
    ├── matched_rules.json
    └── rule_cache
        └── ...
```

The terms "preprocess" and "instrument" represent files generated during two different stages. Please refer to [this document](how-it-works.md) for information about the two stages. For example, `instrument/grpc/clientconn.go` indicates the `clientconn.go` file after code injection. `matched_rules.json` contains the matched rules, and nearly all important files relevant to debugging will be retained in this directory.

## 3. Use delve to debug binary

No optimization will be taken with the `-debug` option during the hybrid compilation. Users can
use [delve](https://github.com/go-delve/delve) to debug the binary file easily.