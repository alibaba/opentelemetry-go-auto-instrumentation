## How to debug

`opentelemetry-go-auto-instrumentation` provides some convenient ways for users to debug the instrumented program.

## 1. Do the instrumentation with -debug options

```bash
# go build
$ ./otel-go-auto-instrumentation -debug
```

After alerting the `-debug` option, `otel-go-auto-instrumentation` will build an unoptimized binary file with more
debugging information.

## 2. Check `.otel-build` directory

After doing the instrumentation, users can check the debugging files in `instrument` directory and `preprocess`
directory. The debugging files
in `preprocess` directory show all the used rule files that are matched and also the log for `go build`. The debugging files in `instrument` directory
show code of the instrumented plugins.

```shell
.otel-build
|-instrument
|--debug_fn_xx.go
|--debug_xx_trampoline.go
|-preprocess
|--backups
|--embededfs
|--debug_otel_rule_xx.go
|--debug_otel_setup_inst.go
|--debug_otel_setup_sdk.go
|--dry_run.log
|--used_rules.json
```

## 3. Use delve to debug binary

No optimization will be taken with the `-debug` option during the hybrid compilation. Users can
use [delve](https://github.com/go-delve/delve) to debug the binary file easily.