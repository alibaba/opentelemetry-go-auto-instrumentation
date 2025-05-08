# 如何添加新规则

## 1. 创建规则
我们通过一个示例演示如何为一个新包注入代码。

以下是如何将日志注入到 `os.Setenv()` 函数中，以跟踪键和值的使用。

首先，在一个新目录中创建并初始化一个 Go 模块：
```console
$ mkdir mysetenv && cd mysetenv
$ go mod init mysetenv
```

接下来，在 `mysetenv` 目录中创建 `hook.go` 文件：
```go
package mysetenv

import (
	"fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func onEnterSetenv(call api.CallContext, key, value string) {
	fmt.Printf("Setting environment variable %s to %s", key, value)
}
```

运行 `go mod tidy` 来下载依赖项。这将设置 hook 代码。

## 2. 注册规则

创建一个 `rule.json` 文件，指定目标函数和 hook 代码：
```json
[
  {
    "ImportPath": "os",
    "Function": "Setenv",
    "OnEnter": "onEnterSetenv",
    "Path": "/path/to/mysetenv"
  }
]
```

- `ImportPath`: 包含要注入函数的包的导入路径
- `Function`: 要注入的函数。
- `OnEnter`: hook 代码
- `Path`: 包含 hook 代码的目录。

还可以指定其他字段，如 `ReceiverType`、`OnExit` 和 `Order`。有关详细信息，请参阅 [文档](rule_def.md)。

## 3. 验证规则

使用一个简单的程序测试该规则：

```console
$ mkdir setenv-demo && cd setenv-demo
$ go mod init setenv-demo
$ cat <<EOF > main.go
package main
import "os"
func main() {
    os.Setenv("hello", "world")
}
EOF
$ ~/otel set -rule=rule.json
$ ~/otel build main.go
$ ./main
Setting environment variable hello to world%
```

输出确认代码注入成功，因为当调用 `os.Setenv()` 时，消息会出现。

此示例未涵盖许多高级功能，您可以参考 `pkg/rules` 目录中的现有规则获取更多信息。