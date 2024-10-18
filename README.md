# OpenTelemetry Go Auto Instrumentation

<img src="docs/logo.png" height="150" align="right" alt="logo">

[![](https://shields.io/badge/Docs-English-blue?logo=Read%20The%20Docs)](./docs)
[![](https://shields.io/badge/Readme-中文-blue?logo=Read%20The%20Docs)](./docs/README_CN.md)

This project provides an automatic solution for Golang applications that want to
leverage OpenTelemetry to enable effective observability. No code changes are
required in the target application, and the instrumentation is done at compile
time. Simply replace `go build` with `otelbuild` to get started.

# Installation

#### Install via Bash
For Linux and MacOS users, use bash to install:
```bash
sudo curl -fsSL https://raw.githubusercontent.com/alibaba/opentelemetry-go-auto-instrumentation/main/install.sh | sudo bash
```

#### Precompiled Binary

Please download the latest precompiled release version from
the [Release](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases)
page.

#### Build From Source

You can also clone the repository and build the project by running the following commands:

```bash
$ make build
```

# Getting Started

Replace `go build` with the following command to build your project:

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
$ ./otelbuild -help        # print help doc
$ ./otelbuild -debuglog    # print log to file
$ ./otelbuild -verbose -- -gcflags="-m" cmd/app # print verbose log
```

You can also explore [**these examples**](./example/) to get hands-on experience.

Also there are several [**documents**](./docs) that you may find useful

> [!NOTE]
> If you find any compilation failures during the process, it's likely a bug.
> Please feel free to file a bug
> at [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
> to help us enhance this project.

# Supported Libraries

| Plugin Name  | Repository Url                             | Min Supported Version | Max Supported Version |
|--------------|--------------------------------------------|-----------------------|-----------------------|
| database/sql | https://pkg.go.dev/database/sql            | -                     | -                     |
| echo         | https://github.com/labstack/echo           | v4.0.0                | v4.12.0               |
| fasthttp     | https://github.com/valyala/fasthttp        | v1.45.0               | v1.55.0               |
| gin          | https://github.com/gin-gonic/gin           | v1.7.0                | v1.10.0               |
| go-redis     | https://github.com/redis/go-redis          | v9.0.5                | v9.5.1                |
| gorm         | https://github.com/go-gorm/gorm            | v1.22.0               | v1.25.9               |
| grpc         | https://google.golang.org/grpc             | v1.44.0               | v1.67.0               |
| hertz        | https://github.com/cloudwego/hertz         | v0.8.0                | v0.9.2                |
| kratos       | https://github.com/go-kratos/kratos        | v2.6.3                | v2.8.1                |
| log          | https://pkg.go.dev/log                     | -                     | -                     |
| logrus       | https://github.com/sirupsen/logrus         | v1.5.0                | v1.9.3                |
| mongodb      | https://github.com/mongodb/mongo-go-driver | v1.11.1               | v1.15.2               |
| mux          | https://github.com/gorilla/mux             | v1.3.0                | v1.8.1                |
| net/http     | https://pkg.go.dev/net/http                | -                     | -                     |
| slog         | https://pkg.go.dev/log/slog                | -                     | -                     |
| zap          | https://github.com/uber-go/zap             | v1.20.0               | v1.27.0               |

We are progressively open-sourcing the libraries we have supported, and your contributions are very welcome 💖!

Please refer to [this document](./docs/how-to-add-a-new-rule.md) for guidance on how to write instrumentation
code for new frameworks.

# Community

We are looking forward to your feedback and suggestions. You can join
our [DingTalk group](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? )
to engage with us.

<img src="docs/dingtalk.png" height="200">
