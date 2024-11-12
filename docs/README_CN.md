![](anim-logo.svg)

[![](https://shields.io/badge/Docs-English-blue?logo=Read%20The%20Docs)](./docs)
[![](https://shields.io/badge/Readme-中文-blue?logo=Read%20The%20Docs)](./README_CN.md)
[![codecov](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation/branch/main/graph/badge.svg)](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation)

该项目为希望利用 OpenTelemetry 的 Golang 应用程序提供了一个自动解决方案。
利用 OpenTelemetry 实现有效可观察性的 Golang 应用程序提供自动解决方案。目标应用程序无需更改代码
在编译时完成。
时完成。只需在 `go build` 中添加 `otelbuild` 前缀即可开始 :rocket：

# 安装

### 通过 Bash 安装
对于 **Linux 和 MacOS** 用户，运行以下命令即可安装该工具
```bash
$ sudo curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/opentelemetry-go-auto-instrumentation@main/install.sh | sudo bash
```
默认情况下，它将安装在 `/usr/local/bin/otelbuild`中。

### 预编译二进制文件

请从
Release](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases)
页面下载最新的预编译版本。

### 从源代码编译

通过运行以下命令查看源代码并构建工具：

```bash
$ make build
```

### 开始

在 `go build` 中添加 `otelbuild` 前缀，以构建项目：

```bash
$ otelbuild go build
$ otelbuild go build -o app cmd/app
$ otelbuild go build -gcflags="-m" cmd/app
```
工具本身的参数应放在 `go build` 之前：

```bash
$ otelbuild -help # 打印帮助文档
$ otelbuild -debug go build # 启用调试模式
$ otelbuild -verbose go build # 打印详细日志
$ otelbuild -rule=custom.json go build # 使用自定义规则
```

您还可以探索 [**这些示例**](./example/) 以获得实践经验。

此外，还有一些 [**文档**](./docs)，您可能会发现它们对了解项目或为项目做出贡献非常有用。

> 注意
> 如果你发现任何编译失败，而 `go build` 却能正常工作，这很可能是一个 bug。
> 请随时在
> 请随时在 [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
> 以帮助我们改进本项目。

# 支持的库

| 插件名称 | 存储库网址 | 最低支持版本 | 最高支持版本
| ------------ | ------------------------------------------ | --------------------- | --------------------- |
| 数据库/sql | https://pkg.go.dev/database/sql | - | - |
| echo | https://github.com/labstack/echo | v4.0.0 | v4.12.0 | v4.12.0
| fasthttp | https://github.com/valyala/fasthttp | v1.45.0 | v1.57.0 |
| gin | https://github.com/gin-gonic/gin | v1.7.0 | v1.10.0 | v4.0.0 | v4.12.0 | fasthttp
| go-redis | https://github.com/redis/go-redis | v9.0.5 | v9.5.1 |
| Gorm | https://github.com/go-gorm/gorm | v1.22.0 | v1.25.9 |
|grpc | https://google.golang.org/grpc | v1.44.0 | v1.67.0 |
|hertz | https://github.com/cloudwego/hertz | v0.8.0 | v0.9.2 |
|kratos | https://github.com/go-kratos/kratos | v2.6.3 | v2.8.2 |
| log | https://pkg.go.dev/log | - | - |
| logrus | https://github.com/sirupsen/logrus | v1.5.0 | v1.9.3 | | mongodb
| mongodb | https://github.com/mongodb/mongo-go-driver | v1.11.1 | v1.15.2 |
| mux | https://github.com/gorilla/mux | v1.3.0 | v1.8.1 |
| net/http | https://pkg.go.dev/net/http | - | - |
| slog | https://pkg.go.dev/log/slog | - | - |
| zap | https://github.com/uber-go/zap | v1.20.0 | v1.27.0 |

我们正在逐步开源我们支持的库，非常欢迎您的贡献💖！

> 重要事项
> 您期望的框架不在列表中？别担心，你可以轻松地将代码注入到任何官方不支持的框架/库中。
>
> 请参考 [this document](./how-to-add-a-new-rule.md) 开始使用。

# 社区

我们期待您的反馈和建议。您可以加入我们的 [DingTalk 群组](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? )
与我们交流。

<img src="dingtalk.png" height="200">