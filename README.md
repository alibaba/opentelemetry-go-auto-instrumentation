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

Lastly, we've provided some examples in `example` direcory, you can try the project through those examples.

### Compatibility

OpenTelemetry-Go Contrib ensures compatibility with the current supported
versions of
the [Go language](https://golang.org/doc/devel/release#policy):

> Each major Go release is supported until there are two newer major releases.
> For example, Go 1.5 was supported until the Go 1.7 release, and Go 1.6 was supported until the Go 1.8 release.

For versions of Go that are no longer supported upstream, opentelemetry-go-contrib will
stop ensuring compatibility with these versions in the following manner:

- A minor release of opentelemetry-go-contrib will be made to add support for the new
  supported release of Go.
- The following minor release of opentelemetry-go-contrib will remove compatibility
  testing for the oldest (now archived upstream) version of Go. This, and
  future, releases of opentelemetry-go-contrib may include features only supported by
  the currently supported versions of Go.

This project is tested on the following systems.

| OS       | Go Version | Architecture |
| -------- | ---------- | ------------ |
| Ubuntu   | 1.22       | amd64        |
| Ubuntu   | 1.21       | amd64        |
| Ubuntu   | 1.22       | 386          |
| Ubuntu   | 1.21       | 386          |
| macOS 13 | 1.22       | amd64        |
| macOS 13 | 1.21       | amd64        |
| macOS    | 1.22       | arm64        |
| macOS    | 1.21       | arm64        |
| Windows  | 1.22       | amd64        |
| Windows  | 1.21       | amd64        |
| Windows  | 1.22       | 386          |
| Windows  | 1.21       | 386          |

While this project should work for other systems, no compatibility guarantees
are made for those systems currently.

# Community

We are looking forward to your feedback and suggestions. Please feel free to join our [DingTalk group](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? ) to communicate with us.

<img src="docs/dingtalk.png" height="200">

Also there are several documents that you may find useful:

- [How it works](./docs/how-it-works.md)
- [How to add a new rule](./docs/how-to-add-a-new-rule.md)
- [How to debug](./docs/how-to-debug.md)
- [How to write tests for plugins](./docs/how-to-write-tests-for-plugins.md)
- [Supported Libraries](./docs/supported-libraries.md)
- [Discussion on this topic at OpenTelemetry community](https://github.com/open-telemetry/community/issues/1961)
- [面向OpenTelemetry的Golang应用无侵入插桩技术](https://mp.weixin.qq.com/s/FKCwzRB5Ujhe1stOH2ibXg)
