# OpenTelemetry Go Auto Instrumentation

# Introduction

This project provides an automatic solution for Golang applications that want to
leverage OpenTelemetry to enable effective observability.

# How to Build

Run the following command to build `otel-go-auto-instrumentation`:

```bash
$ make build
```

For all supported platforms:

```bash
$ make all
```

To run all tests:

```bash
$ make test
```

# How to Use

Replace `go build` with the following command to build you project:

```bash
# go build
$ ./otel-go-auto-instrumentation
```

The arguments for `go build` should be placed after the `--` delimiter:

```bash
# go build -gcflags="-m" cmd/app
$ ./otel-go-auto-instrumentation -- -gcflags="-m" cmd/app
```

The arguments for the tool itself should be placed before the `--` delimiter:

```bash
$ ./otel-go-auto-instrumentation -help # print help doc
$ ./otel-go-auto-instrumentation -debuglog # print log to file
$ ./otel-go-auto-instrumentation -verbose -- -gcflags="-m" cmd/app # print verbose log
```

If you find any failures during the process, it's likely a bug.
Please feel free to file a bug
at [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
to help us enhance this project.

# Community

We are looking forward to your feedback and suggestions. Please feel free to join our [DingTalk group](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? ) to communicate with us.

<img src="docs/dingtalk.png" height="200">

Also there are several documents that you may find useful:

- [How it works](./docs/how-it-works.md)
- [How to add a new rule](./docs/how-to-add-a-new-rule.md)
- [Discussion on this topic at OpenTelemetry community](https://github.com/open-telemetry/community/issues/1961)