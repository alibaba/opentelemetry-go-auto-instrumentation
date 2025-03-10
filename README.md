![](docs/anim-logo.svg)

[![](https://shields.io/badge/Docs-English-blue?logo=Read%20The%20Docs)](./docs/README.md)
[![](https://shields.io/badge/Readme-中文-blue?logo=Read%20The%20Docs)](./docs/README_CN.md)
[![codecov](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation/branch/main/graph/badge.svg)](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation)

This project provides an automatic solution for Golang applications that want to
leverage OpenTelemetry to enable effective observability. No code changes are
required in the target application, the instrumentation is done at compile
time. Simply adding `otel` prefix to `go build` to get started :rocket:

# Installation

### Install via Bash
For **Linux and MacOS** users, install the tool by running the following command
```bash
$ sudo curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/opentelemetry-go-auto-instrumentation@main/install.sh | sudo bash
```
It will be installed in `/usr/local/bin/otel` by default.

### Precompiled Binary

Please download the latest precompiled release version from
the [Release](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases)
page.

### Build From Source

Checkout source code and build the tool by running one of following commands:

```bash
$ make         # build only
$ make install # build and install
```

# Getting Started

Check the version by running:
```bash
$ otel version
```

The configuration for the tool can be set by the following command:

```bash
$ otel set -verbose                          # print verbose logs
$ otel set -log=/path/to/file.log            # set log file
$ otel set -debug                            # enable debug mode
$ otel set -debug -verbose -rule=custom.json # set multiple configs
$ otel set -disabledefault -rule=custom.json # disable default rules, use custom rules only
$ otel set -rule=custom.json                 # use default and custom rules
$ otel set -rule=a.json,b.json               # use default, a and b rules
```

Once all set up, add `otel` prefix to `go build` to build your project:

```bash
$ otel go build
$ otel go build -o app cmd/app
$ otel go build -gcflags="-m" cmd/app
```

The detailed usage of otel tool can be found in [**Usage**](./docs/usage.md).

> [!NOTE]
> If you find any compilation failures while `go build` works, it's likely a bug.
> Please feel free to file a bug
> at [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
> to help us enhance this project.

# Examples

You can also explore [**these examples**](./example/) to get hands-on experience. They are designed to construct a full picture of how to use the tool in different scenarios.

Also there are several [**documents**](./docs) that you may find useful for either understanding the project or contributing to it.

# Supported Libraries

| Plugin Name   | Repository Url                                 | Min Supported Version | Max Supported Version |
|---------------| ---------------------------------------------- |-----------------------|-----------------------|
| database/sql  | https://pkg.go.dev/database/sql                | -                     | -                     |
| echo          | https://github.com/labstack/echo               | v4.0.0                | v4.12.0               |
| elasticsearch | https://github.com/elastic/go-elasticsearch    | v8.4.0                | v8.15.0               |
| fasthttp      | https://github.com/valyala/fasthttp            | v1.45.0               | v1.57.0               |
| fiber         | https://github.com/gofiber/fiber               | v2.43.0               | v2.52.5               |
| gin           | https://github.com/gin-gonic/gin               | v1.7.0                | v1.10.0               |
| gorestful     | https://github.com/emicklei/go-restful         | v3.7.0                | v3.12.1               |
| go-redis      | https://github.com/redis/go-redis              | v9.0.5                | v9.5.1                |
| go-redis v8   | https://github.com/redis/go-redis              | v8.11.0               | v8.11.5               |
| gorm          | https://github.com/go-gorm/gorm                | v1.22.0               | v1.25.9               |
| grpc          | https://google.golang.org/grpc                 | v1.44.0               | v1.68.2               |
| hertz         | https://github.com/cloudwego/hertz             | v0.8.0                | v0.9.2                |
| kitex         | https://github.com/cloudwego/kitex             | v0.5.1                | v0.11.3               |
| kratos        | https://github.com/go-kratos/kratos            | v2.6.3                | v2.8.2                |
| log           | https://pkg.go.dev/log                         | -                     | -                     |
| logrus        | https://github.com/sirupsen/logrus             | v1.5.0                | v1.9.3                |
| mongodb       | https://github.com/mongodb/mongo-go-driver     | v1.11.1               | v1.15.2               |
| mux           | https://github.com/gorilla/mux                 | v1.3.0                | v1.8.1                |
| nacos         | https://github.com/nacos-group/nacos-sdk-go/v2 | v2.0.0                | v2.2.7                |
| net/http      | https://pkg.go.dev/net/http                    | -                     | -                     |
| redigo        | https://github.com/gomodule/redigo             | v1.9.0                | v1.9.2                |
| slog          | https://pkg.go.dev/log/slog                    | -                     | -                     |
| zap           | https://github.com/uber-go/zap                 | v1.20.0               | v1.27.0               |
| zerolog       | https://github.com/rs/zerolog                  | v1.10.0               | v1.33.0               |
| iris          | https://github.com/kataras/iris                | v12.2.0               | v12.2.11              |
| gomicro       | https://github.com/micro/go-micro              | v5.0.0                | v5.3.0                |
| trpc-go       | https://github.com/trpc-group/trpc-go          | v1.0.0                | v1.0.3                |

We are progressively open-sourcing the libraries we have supported, and your contributions are very welcome 💖!

> [!IMPORTANT]
> The framework you expected is not in the list? Don't worry, you can easily inject your code into any frameworks/libraries that are not officially supported.
>
> Please refer to [this document](./docs/how-to-add-a-new-rule.md) to get started.

# Community

We are looking forward to your feedback and suggestions. You can join
our [DingTalk group](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? )
to engage with us.

<img src="docs/dingtalk.png" height="200">

# Adopters

These are only part of the companies using this project, for reference only. If you are using this project, please [add your company here](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues/225) to tell us your scenario to make this project better.

- <img src="./docs/alibaba.png" width="80">
- <img src="./docs/aliyun.png" width="100">