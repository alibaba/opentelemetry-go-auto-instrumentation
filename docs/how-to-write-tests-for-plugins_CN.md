## 如何为插件编写测试

当你按照 [how-to-add-a-new-rule.md](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/blob/main/docs/how-to-add-a-new-rule.md) 添加了新的自动注入规则后，你需要添加相应的测试用例以验证你的规则是否生效。`opentelemetry-go-auto-instrumentation` 提供了一种便捷的方式来完成规则验证。

## 添加通用插件测试用例

切换目录到 `/test`，并为你要测试的插件创建一个新的目录，例如 `redis`。在 `redis` 目录中，子目录的名称应表示该插件支持的最低版本。如果你要为 Redis 插件添加测试，则需要执行以下步骤：

### 1. 添加你的规则支持的最低版本依赖

例如，如果你添加的规则支持从 `v9.0.5` 起的 redis 版本，那么你应首先验证最低版本，即 `v9.0.5`。你可以创建一个名为 `v9.0.5` 的子目录，并添加如下 `go.mod` 文件：
```
module redis/v9.0.5

go 1.22

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../opentelemetry-go-auto-instrumentation

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../../opentelemetry-go-auto-instrumentation/test/verifier

require (
	// import this dependency to use verifier
    github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.0.5
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/sdk v1.30.0
)
```
### 2. 基于插件编写业务逻辑

接下来，你需要基于 `redis` 插件编写一些业务逻辑代码。例如，`test_executing_commands.go` 实现了 Redis 中的基本 `get` 和 `set` 操作。你的测试应尽可能覆盖该插件的各种使用场景。

### 3. 编写验证代码

如果你编写的业务逻辑代码与注入规则匹配，则会生成一些遥测数据（如 span）。例如，`test_executing_commands.go` 应生成两个 span：一个表示 `set` 操作，另一个表示 `get` 操作。你应使用 `verifier` 来验证这些数据的正确性：
```go
import "github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"

verifier.WaitAndAssertTraces(func (stubs []tracetest.SpanStubs) {
	verifier.VerifyDbAttributes(stubs[0][0], "set", "", "redis", "", "localhost", "set a b ex 5 ", "set", "")
	verifier.VerifyDbAttributes(stubs[1][0], "get", "", "redis", "", "localhost", "get a ", "get", "")
})
```

`verifier` 的 `WaitAndAssertTraces` 方法接受一个回调函数，该函数会提供所有生成的 trace。你应在该回调中验证每个 trace 中所有 span 的属性、父级上下文以及其他关键信息。

如果你还想验证指标（metrics）数据，也可以使用 `verifier`，如下所示：
```go
	verifier.WaitAndAssertMetrics(map[string]func(metricdata.ResourceMetrics) {
		"http.server.request.duration": func(mrs metricdata.ResourceMetrics) {
		if len(mrs.ScopeMetrics) <= 0 {
			panic("No http.server.request.duration metrics received!")
		}
		point := mrs.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
		if point.DataPoints[0].Count != 1 {
			panic("http.server.request.duration metrics count is not 1")
		}
		verifier.VerifyHttpServerMetricsAttributes(point.DataPoints[0].Attributes.ToSlice(), "GET", "/a", "", "http", "1.1", "http", 200)
		},
		"http.client.request.duration": func(mrs metricdata.ResourceMetrics) {
		if len(mrs.ScopeMetrics) <= 0 {
			panic("No http.client.request.duration metrics received!")
		}
		point := mrs.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
		if point.DataPoints[0].Count != 1 {
			panic("http.client.request.duration metrics count is not 1")
		}
		verifier.VerifyHttpClientMetricsAttributes(point.DataPoints[0].Attributes.ToSlice(), "GET", "127.0.0.1:"+strconv.Itoa(port), "", "http", "1.1", port, 200)
       },
	})
```
用户需要使用 `verifier` 中的 `WaitAndAssertMetrics` 方法来验证指标数据的正确性。`WaitAndAssertMetrics` 接收一个映射，映射的键是指标的名称，值是该指标数据的验证函数。用户可以在回调函数中编写自己的验证逻辑。

### 4. 注册测试

最后，你需要注册测试。你应在 `test` 目录中编写一个 `_tests.go` 文件来进行注册：

```go
const redis_dependency_name = "github.com/redis/go-redis/v9"
const redis_module_name = "redis"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("redis-9.0.5-executing-commands-test", redis_module_name, "v9.0.5", "v9.5.1", "1.18", "", TestExecutingCommands)
}

func TestExecutingCommands(t *testing.T, env ...string) {
	redisC, redisPort := initRedisContainer()
	defer clearRedisContainer(redisC)
	UseApp("redis/v9.0.5")
	RunGoBuild(t, "go", "build", "test_executing_commands.go")
	env = append(env, "REDIS_PORT="+redisPort.Port())
	RunApp(t, "test_executing_commands", env...)
}

```

在 `init` 函数中，你需要使用 `NewGeneralTestCase` 来包装你的测试用例，该函数接收以下参数：

- testName：测试用例的名称。
- moduleName：`test` 目录中子目录的名称。
- minVersion：插件支持的最低版本。
- maxVersion：插件支持的最高版本。
- minGoVersion：插件支持的最低 Go 版本。
- maxGoVersion：插件支持的最高 Go 版本。
- testFunc：要执行的测试函数。

你应该使用 `opentelemetry-go-auto-instrumentation` 来构建测试用例，使其能够生成遥测数据。首先，你应该调用 `UseApp` 方法切换到测试用例所在的目录，然后调用 `RunGoBuild` 进行混合编译。最后，使用 `RunApp` 运行已注入的测试用例二进制文件，以验证遥测数据。
```go
func TestExecutingUnsupportedCommands(t *testing.T, env ...string) {
	redisC, redisPort := initRedisContainer()
	defer clearRedisContainer(redisC)
	UseApp("redis/v9.0.5")
	RunGoBuild(t, "go", "build", "test_executing_unsupported_commands.go")
	env = append(env, "REDIS_PORT="+redisPort.Port())
	RunApp(t, "test_executing_unsupported_commands", env...)
}
```

## 添加禁用检查用例

禁用检查灵感来源于 [safety-mechanisms.md](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/safety-mechanisms.md)。
因为运行每个版本的常规插件测试会非常耗时，所以 `opentelemetry-go-auto-instrumentation` 会随机选择一个版本进行混合编译，
以验证不同版本之间的 API 兼容性。如果禁用检查发现某些版本的 API 已经发生了变化，社区将创建新的规则来适配它。

用户可以通过调用 `NewMuzzleTestCase` 来添加禁用检查用例，`NewMuzzleTestCase` 接受的参数几乎与 `NewGeneralTestCase` 相同。
你需要额外指定插件的依赖名称和需要进行禁用检查的类列表。

## 添加最新版本深度检查用例

仪器测试通常会针对我们支持的最低版本库进行，以确保用户在使用旧依赖版本时的基准测试。
由于代理的性质以及我们在私有 API 上进行插桩的位置，代理可能会在新发布的库版本上失败。为了避免这种情况，我们还会在夜间构建过程中，针对从远程获取的库的最新版本运行插桩测试。
如果库的新版本无法与代理兼容，我们可以通过这次构建发现问题，并在下一个版本的代理中进行修复。

用户可以通过调用 `NewLatestDepthTestCase` 来添加最新版本深度检查用例，`NewLatestDepthTestCase` 接受的参数几乎与 `NewGeneralTestCase` 相同。
你需要额外指定插件的依赖名称以及需要进行最新版本深度检查的类列表。