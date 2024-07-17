# OpenTelemetry Go Auto Instrumentation

## Introduction

This project provides an automatic solution for Golang applications that want to
leverage OpenTelemetry to enable effective observability.

## Build
To build local development version, run below command:

```bash
$ make build
```

To build release version for all supported os+arch matrix, run below command:

```bash
$ make all
```

## Usage
If you previously used `go build` to build your project, you can simply use

```bash
$ ./otel-go-auto-instrumentation
```

If you previously used `go build` with arguments, e.g.  `go build -gcflags="-m" cmd/app`, it becomes

```bash
$ ./otel-go-auto-instrumentation -- -gcflags="-m" cmd/app
```

Where all arguments that you would pass to `go build` are now passed after the `--` delimiter.

Additionally, there are some options for specific purposes, which need to be passed
before the -- delimiter, as shown below:

```bash
$ otel-go-auto-instrumentation -help
$ otel-go-auto-instrumentation -debuglog # print log to file
$ otel-go-auto-instrumentation -verbose -- -gcflags="-m" cmd/app # print verbose log
```

As a drop-in replacement for `go build`, you can replace `go build` in your project's
Makefile/build script/build command with `otel-go-auto-instrumentation --` and
then compile the project as usual. If you find that the build fails after
the replacement, it is likely a bug, please file a bug in [Issue](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues).
