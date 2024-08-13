# How to add a new rule

## 1. Create a new rule
We demonstrate how to inject code for a new package through an example. Here we choose `os.Setenv()` and inject logging code at the beginning of its function to record the key and value whenever user calls it.

First, we create the following file structure under the `pkg/rules` directory:
```
pkg/rules/os
├── rule.go
└── setup.go
```
The `os` directory is the package name, and `rule.go` is the rule file where we define the rule. `setup.go` is the setup file where we register the rule.

```go
package os

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	api.NewRule("os", "Setenv", "", "onEnterSetenv", "").
		Register()
}
```
In the `rule.go` file, we define a new rule for the `os` package. The `NewRule` function takes five arguments: the package name, the function name, the function signature, the onEnter function name, and the onExit function name. The `Register` function registers the rule.

`setup.go` contains the definition of `onEnterSetenv`:

```go
//go:build ignore
package os

import (
	"fmt"
	"os"
)

func onEnterSetenv(call os.CallContext, key, value string) {
	fmt.Printf("Setting environment variable %s to %s", key, value)
}
```

The `onEnterSetenv` function is the probe code that will be injected at the beginning of the `os.Setenv()` function. It takes the `os.CallContext` as the first argument and the key and value as the second and third arguments, respectively.

## 2. Register the rule
After completing this process, we introduce our newly created package in `rule_enabler.go` to make the rule take effect:

```go
package main
import (
    ...
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/os"
)
```
Dont forget recomplie the tool after adding the new rule, because currently `rules` is coded within the tool.

## 3. Verify the rule
To make sure the rule is working as expected, we can write a simple demo program to test it:

```bash
$ mkdir setenv
$ go mod init setenv
$ cat <<EOF > main.go
package main
import "os"
func main(){ os.Setenv("hello", "world") }
EOF
$ ~/opentelemetry-go-auto-instrumentation/otel-go-auto-instrumentation -- main.go
$ ./main
Setting environment variable hello to world%
```
The output shows that the rule is working as expected. The `Setting environment variable hello to world` message is printed when the `os.Setenv()` function is called, that is, the probe code is injected successfully.

There are many advanced features that is not covered in this example, such as the `InstStructRule`, `InstFileRule` and APIs of `CallContext`. You can refer to the existing rules in the `pkg/rules` directory for more information.