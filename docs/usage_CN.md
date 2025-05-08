# `Otel` 使用指南

## 介绍
本指南提供了详细的配置和使用 `otel` 工具的概述。该工具允许你设置各种配置选项，构建项目，并自定义工作流程以获得最佳性能。

## 配置
配置工具的主要方法是通过 `otel set` 命令。此命令允许你指定针对需求量身定制的各种设置：

- **日志**：设置一个自定义的日志文件，用于存储工具生成的日志。
```console
  $ otel set -log=/path/to/file.log
```
默认的日志文件是 `.otel-build/preprocess/debug.log`。你可以通过设置 `-log=/dev/stdout` 将日志输出到标准输出。

**详细日志**：启用详细日志记录，以获取工具的详细输出，这对于故障排除和理解工具的处理过程非常有帮助。
```console
  $ otel set -verbose
```

**调试模式**：启用调试模式以收集调试级别的洞察和信息。
```console
  $ otel set -debug
```

**多重配置**：一次设置多个配置。例如，在使用自定义规则文件时，启用调试模式和详细日志模式：
```console
  $ otel set -debug -verbose -rule=custom.json
```

**仅使用自定义规则**：禁用默认规则集，仅应用特定的自定义规则。当你需要为项目定制规则集时，这特别有用。
```console
  $ otel set -disabledefault -rule=custom.json
```

**默认规则和自定义规则的组合**：同时使用默认规则和自定义规则，以提供全面的配置：
```console
  $ otel set -rule=custom.json
```

**多个规则文件**：将多个自定义规则文件与默认规则结合使用，规则文件可以通过逗号分隔的列表指定：
```console
  $ otel set -rule=a.json,b.json
```

## 使用环境变量
除了使用 `otel set` 命令外，配置还可以通过环境变量进行覆盖。例如，`OTELTOOL_DEBUG` 环境变量允许您临时强制工具进入调试模式，这种方法对于一次性配置非常有效，而无需更改永久设置。
```console
$ export OTELTOOL_DEBUG=true
$ export OTELTOOL_VERBOSE=true
```

环境变量的名称对应于 `otel set` 命令中可用的配置选项，并以 `OTELTOOL_` 为前缀。

环境变量完整列表：

- `OTELTOOL_DEBUG`：启用调试模式。
- `OTELTOOL_VERBOSE`：启用详细日志记录。
- `OTELTOOL_RULE_JSON_FILES`：指定自定义规则文件。
- `OTELTOOL_DISABLE_DEFAULT`：禁用默认规则。
  这种方法为测试更改和实验配置提供了灵活性，您可以在不永久更改现有设置的情况下进行调整。
## 构建项目
一旦配置到位，您可以使用带有 `otel` 前缀的命令来构建项目。这将工具的配置直接集成到构建过程中：

标准构建：使用默认设置构建您的项目。
```console
  $ otel go build
```

输出到特定位置：构建您的项目并指定输出位置。
```console
  $ otel go build -o app cmd/app
```

传递编译器标志：使用编译器标志进行更定制化的构建。
```console
  $ otel go build -gcflags="-m" cmd/app
```
无论您的项目有多复杂，otel 工具都能简化这个过程，通过自动对您的代码进行仪表化，以实现有效的可观察性，唯一的要求是将 `otel` 前缀添加到您的构建命令中。