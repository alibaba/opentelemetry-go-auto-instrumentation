# OpenTelemetry Go Auto Instrumentation

<img src="docs/logo.png" height="150" align="right" alt="logo">

[![](https://shields.io/badge/Docs-English-blue?logo=Read%20The%20Docs)](./docs)
[![](https://shields.io/badge/Readme-中文-blue?logo=Read%20The%20Docs)](./docs/README_CN.md)

为Golang应用程序提供一种面向OpenTelemetry的可观测解决方案。目标应用程序不需要进行代码更改，插桩工作在编译时完成。只需将`go build`替换为`otelbuild`即可开始。

# 构建

运行以下命令来构建`otelbuild`:

```bash
$ make build
```

要为发布目的准备所有支持平台的二进制文件，请使用:

```bash
$ make all
```

运行所有测试:

```bash
$ make test
```

# 使用

使用以下命令替换`go build`来构建你的项目

```bash
# go build
$ ./otelbuild
```

`go build`的参数应放在`--`分隔符之后：

```bash
# go build -gcflags="-m" cmd/app
$ ./otelbuild -- -gcflags="-m" cmd/app
```

工具本身的参数应放在`--`分隔符之前:

```bash
$ ./otelbuild -help        # print help doc
$ ./otelbuild -debuglog    # print log to file
$ ./otelbuild -verbose -- -gcflags="-m" cmd/app # print verbose log
```

你也可以浏览[这些示例](./example/)快速上手.

> [!NOTE]
> 如果过程中发现任何编译失败的情况，很可能是一个bug
> 请在[GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
> 中提交 bug，帮助我们改进这个项目

# 支持的库

| 插件名字  | 仓库地址                              | 最低版本 | 最高版本 |
|--------------|---------------------------------------------|-----------------------|-----------------------|
| database/sql | https://pkg.go.dev/database/sql             | -                     | -                     |
| echo         | https://github.com/labstack/echo            | v4.0.0                | v4.12.0               |
| gin          | https://github.com/gin-gonic/gin            | v1.7.0                | v1.10.0               |
| go-redis     | https://github.com/redis/go-redis           | v9.0.5                | v9.5.1                |
| gorm         | https://github.com/go-gorm/gorm             | v1.22.0               | v1.25.9               |
| logrus       | https://github.com/sirupsen/logrus          | v1.5.0                | v1.9.3                |
| mongodb      | https://github.com/mongodb/mongo-go-driver  | v1.11.1               | v1.15.2               |
| mux          | https://github.com/gorilla/mux              | v1.3.0                | v1.8.1                |
| net/http     | https://pkg.go.dev/net/http                 | -                     | -                     |
| zap          | https://github.com/uber-go/zap              | v1.20.0               | v1.27.0               |


我们正在逐步开源我们支持的库，也非常欢迎你的贡献。有关如何为新框架编写插桩代码的指南，请参阅[这个文档](./docs/how-to-add-a-new-rule.md)

# Community

我们期待你的反馈和建议。请加入我们的[钉钉群](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? )
与我们交流。

<img src="docs/dingtalk.png" height="200">

此外，还有以下一些你可能会觉得有用的文档：

- [如何添加新规则](./docs/how-to-add-a-new-rule.md)
- [如何为插件编写测试](./docs/how-to-write-tests-for-plugins.md)
- [兼容性](./docs/compatibility.md)
- [工作原理](./docs/how-it-works.md)
- [调试指南](./docs/how-to-debug.md)
- [上下文传播](./docs/context-propagation.md)
- [支持的库](./docs/supported-libraries.md)
- [Benchmark](./example/benchmark/benchmark.md)
- [在OpenTelemetry社区讨论该项目的捐赠](https://github.com/open-telemetry/community/issues/1961)
- [面向OpenTelemetry的Golang应用无侵入插桩技术](https://mp.weixin.qq.com/s/FKCwzRB5Ujhe1stOH2ibXg)
