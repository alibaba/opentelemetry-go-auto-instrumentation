# How to add a new rule

## 1. Create a new rule
We demonstrate how to inject code for a new package through an example.
Here we choose `os.Setenv()` and inject logging code at the beginning of its 
function to record the key and value whenever user calls it.

First, we create the following a `mysetenv` directory and initialize it as a Go module:

```bash
$ mkdir mysetenv && cd mysetenv
$ go mod init mysetenv
```

Then we write the following `hook.go` file in the `mysetenv` directory:

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

Run `go mod tidy` to download the dependencies. 
That's the whole story for the probe code.

## 2. Register the rule
We create a new `rule.json` file and specify which function we want to instrument
and which probe code to inject, etc.:

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

- `ImportPath`: The import path of the package that contains the function to be instrumented.
- `Function`: The name of the function to be instrumented.
- `OnEnter`: The name of the function to be called when the instrumented function is called.
- `Path`: The path to the directory containing the probe code.

There are many other fields that can be specified in the rule, such as 
`OnExit`, `Order`, etc. You can refer to [this document](rule_def.md) for more information.

## 3. Verify the rule
To make sure the rule is working as expected, we can write a simple demo program to test it:

```bash
$ mkdir setenv && cd setenv
$ go mod init setenv
$ cat <<EOF > main.go
package main
import "os"
func main(){ os.Setenv("hello", "world") }
EOF
$ ~/otelbuild -rule=rule.json -- main.go
$ ./main
Setting environment variable hello to world%
```
The output shows that the rule is working as expected. The `Setting environment variable hello to world` message is printed when the `os.Setenv()` function is called, that is, the probe code is injected successfully.

There are many advanced features that is not covered in this example, such as the `InstStructRule`, `InstFileRule` and APIs of `CallContext`. You can refer to the existing rules in the `pkg/rules` directory for more information.
