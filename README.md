# OpenTelemetry Go Auto Instrumentation

<img src="docs/logo.png" height="150" align="right">

# Introduction

This project provides an automatic solution for Golang applications that want to
leverage OpenTelemetry to enable effective observability.

# How to Build

Run the following command to build `otelbuild`:

```bash
$ make build
```

To build all supported platforms for release purpose:

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
$ ./otelbuild
```

The arguments for `go build` should be placed after the `--` delimiter:

```bash
# go build -gcflags="-m" cmd/app
$ ./otelbuild -- -gcflags="-m" cmd/app
```

The arguments for the tool itself should be placed before the `--` delimiter:

```bash
$ ./otelbuild -help # print help doc
$ ./otelbuild -debuglog # print log to file
$ ./otelbuild -verbose -- -gcflags="-m" cmd/app # print verbose log
```

If you find any failures during the process, it's likely a bug.
Please feel free to file a bug
at [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
to help us enhance this project.

Lastly, we've provided some examples in [example](/example/) direcory, you can try the project through those examples.

# Supported libraries

| Plugin Name  | Repository Url                             | Min Supported Version | Max Supported Version |
|--------------|--------------------------------------------|-----------------------|-----------------------|
| net/http     | https://pkg.go.dev/net/http                | -                     | -                     |
| database/sql | https://pkg.go.dev/database/sql            | -                     | -                     |
| go-redis     | https://github.com/redis/go-redis          | v9.0.5                | v9.5.1                |
| mongodb      | https://github.com/mongodb/mongo-go-driver | v1.11.1               | v1.15.2               |

We are gradually open-sourcing the libraries we have supported, and your contributions are warmly welcome.

# Community

We are looking forward to your feedback and suggestions. Please feel free to join our [DingTalk group](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? ) to communicate with us.

<img src="docs/dingtalk.png" height="200">

Also there are several documents that you may find useful:

- [How it works](./docs/how-it-works.md)
- [How to add a new rule](./docs/how-to-add-a-new-rule.md)
- [How to debug](./docs/how-to-debug.md)
- [How to write tests for plugins](./docs/how-to-write-tests-for-plugins.md)
- [Compatibility](./docs/compatibility.md)
- [Context Propagation](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/blob/main/docs/context-propagation.md)
- [Supported Libraries](./docs/supported-libraries.md)
- [Discussion on this topic at OpenTelemetry community](https://github.com/open-telemetry/community/issues/1961)
- [面向OpenTelemetry的Golang应用无侵入插桩技术](https://mp.weixin.qq.com/s/FKCwzRB5Ujhe1stOH2ibXg)
