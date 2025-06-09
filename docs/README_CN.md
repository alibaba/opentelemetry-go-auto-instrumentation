![](anim-logo.svg)

[![](https://shields.io/badge/Docs-English-blue?logo=Read%20The%20Docs)](./README.md) &nbsp;
[![](https://shields.io/badge/Readme-中文-blue?logo=Read%20The%20Docs)](./README_CN.md)  &nbsp;
[![codecov](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation/branch/main/graph/badge.svg)](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation)  &nbsp;
[![](https://shields.io/badge/Aliyun-Commercial-orange?logo=alibabacloud)](https://help.aliyun.com/zh/arms/application-monitoring/getting-started/monitoring-the-golang-applications) &nbsp;
[![](https://img.shields.io/badge/New-Adopter-orange?logo=githubsponsors)](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues/225) &nbsp;

该项目为希望利用 OpenTelemetry 的 Golang 应用程序提供了一个自动解决方案。
利用 OpenTelemetry 实现有效可观察性的 Golang 应用程序提供自动解决方案。目标应用程序无需更改代码
在编译时完成。只需在 `go build` 中添加 `otel` 前缀即可开始 :rocket:

# 安装

### 通过 Bash 安装
对于 **Linux 和 MacOS** 用户，运行以下命令即可安装该工具
```console
$ sudo curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/opentelemetry-go-auto-instrumentation@main/install.sh | sudo bash
```
默认情况下，它将安装在 `/usr/local/bin/otel`中。

### 预编译二进制文件

请从
[Release](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases)
页面下载最新的预编译版本。

### 从源代码编译

通过运行以下命令查看源代码并构建工具：

```console
$ make         # 只构建
$ make install # 构建并安装
```

# 开始

通过运行以下命令检查版本：
```console
$ otel version
```

通过以下命令配置工具参数：
```console
$ otel set -verbose                          # 打印详细日志
$ otel set -debug                            # 启用调试模式
$ otel set -rule=custom.json                 # 同时使用默认和自定义规则
```

通常情况下，你无需设置任何配置。只需在 `go build` 前加上 `otel` 前缀来构建你的项目：

```console
$ otel go build
$ otel go build -o app cmd/app
$ otel go build -gcflags="-m" cmd/app
```

这就是整个过程！该工具将自动使用 OpenTelemetry 对你的代码进行插装，你就可以开始观测你的应用程序了。 : 望远镜:

您可以在 [**使用指南**](./usage.md)中找到 `otel` 工具的详细用法。

> [!NOTE] 
> 如果您发现任何编译失败，而 `go build` 却能正常工作，这很可能是一个 bug。
> 请随时在
> [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
> 提交问题报告以帮助我们改进本项目。

# 示例

您还可以探索 [**这些示例**](../example/) 以获得实践经验。

此外，还有一些 [**文档**](./)，您可能会发现它们对了解项目或为项目做出贡献非常有用。

# 支持的库

| 插件名称          | 存储库网址                                          | 最低支持版本  | 最高支持版本   |
|---------------|------------------------------------------------|---------|----------|
| database/sql  | https://pkg.go.dev/database/sql                | -       | -        |
| echo          | https://github.com/labstack/echo               | v4.0.0  | v4.12.0  |
| elasticsearch | https://github.com/elastic/go-elasticsearch    | v8.4.0  | v8.15.0  |
| fasthttp      | https://github.com/valyala/fasthttp            | v1.45.0 | v1.59.0  |
| fiber         | https://github.com/gofiber/fiber               | v2.43.0 | v2.52.6  |
| gin           | https://github.com/gin-gonic/gin               | v1.7.0  | v1.10.0  |
| go-redis      | https://github.com/redis/go-redis              | v9.0.5  | v9.5.1   |
| go-redis v8   | https://github.com/redis/go-redis              | v8.11.0 | v8.11.5  |
| gomicro       | https://github.com/micro/go-micro              | v5.0.0  | v5.3.0   |
| gorestful     | https://github.com/emicklei/go-restful         | v3.7.0  | v3.12.1  |
| gorm          | https://github.com/go-gorm/gorm                | v1.22.0 | v1.25.9  |
| grpc          | https://google.golang.org/grpc                 | v1.44.0 | v1.71.0  |
| hertz         | https://github.com/cloudwego/hertz             | v0.8.0  | v0.9.2   |
| iris          | https://github.com/kataras/iris                | v12.2.0 | v12.2.11 |
| kitex         | https://github.com/cloudwego/kitex             | v0.5.1  | v0.11.3  |
| kratos        | https://github.com/go-kratos/kratos            | v2.6.3  | v2.8.4   |
| langchaingo   | https://github.com/tmc/langchaingo             | v0.1.13 | v0.1.13  |
| log           | https://pkg.go.dev/log                         | -       | -        |
| logrus        | https://github.com/sirupsen/logrus             | v1.5.0  | v1.9.3   |
| mongodb       | https://github.com/mongodb/mongo-go-driver     | v1.11.1 | v1.15.1  |
| mux           | https://github.com/gorilla/mux                 | v1.3.0  | v1.8.1   |
| nacos         | https://github.com/nacos-group/nacos-sdk-go/v2 | v2.0.0  | v2.2.7   |
| net/http      | https://pkg.go.dev/net/http                    | -       | -        |
| redigo        | https://github.com/gomodule/redigo             | v1.9.0  | v1.9.2   |
| slog          | https://pkg.go.dev/log/slog                    | -       | -        |
| trpc-go       | https://github.com/trpc-group/trpc-go          | v1.0.0  | v1.0.3   |
| zap           | https://github.com/uber-go/zap                 | v1.20.0 | v1.27.0  |
| zerolog       | https://github.com/rs/zerolog                  | v1.10.0 | v1.33.0  |
| rocketmq      | https://github.com/apache/rocketmq-client-go   | v2.0.0  | v2.1.2   |

我们正在逐步开源我们支持的库，非常欢迎您的贡献💖！

> [!IMPORTANT]
> 您期望的框架不在列表中？别担心，您可以轻松地将代码注入到任何官方不支持的框架/库中。
>
> 请参考 [自定义插桩文档](./how-to-add-a-new-rule.md) 开始使用。

# 社区

我们期待您的反馈和建议。您可以加入我们的 [钉钉群组](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? )
与我们交流。

<img src="dingtalk.png" height="200">

我们衷心感谢以下为该项目做出宝贵贡献的贡献者：

<a href="https://github.com/alibaba/opentelemetry-go-auto-instrumentation/graphs/contributors">
  <img alt="contributors" src="https://contrib.rocks/image?repo=alibaba/opentelemetry-go-auto-instrumentation" height="100"/>
</a>

该项目的Star历史如下，它展示了这个项目随着时间推移的发展情况：

[![Star History](https://api.star-history.com/svg?repos=alibaba/opentelemetry-go-auto-instrumentation&type=Date)](https://star-history.com/#alibaba/opentelemetry-go-auto-instrumentation&Date)