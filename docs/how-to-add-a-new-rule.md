# How to add a new rule

## 1. Create the Rule
We demonstrate how to inject code for a new package through an example.

Here's how to inject logging into the os.Setenv() function to track key and value usage.

First, create and initialize a Go module in a new directory:

```bash
$ mkdir mysetenv && cd mysetenv
$ go mod init mysetenv
```

Next, create the hook.go file in the mysetenv directory:

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

Run `go mod tidy` to download the dependencies. This sets up the hook code.

## 2. Register the Rule

Create a `rule.json` file to specify the target function and hook code:

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

- `ImportPath`: The import path of the package that contains the function
- `Function`: The function to instrument.
- `OnEnter`: The hook code
- `Path`: Directory containing the hook code.

Additional fields like `OnExit` and `Order` can also be specified. Refer to [the documentation](rule_def.md) for details.

## 3. Verify the Rule
Test the rule with a simple program:

```bash
$ mkdir setenv-demo && cd setenv-demo
$ go mod init setenv-demo
$ cat <<EOF > main.go
package main
import "os"
func main() {
    os.Setenv("hello", "world")
}
EOF
$ ~/otel -rule=rule.json go build main.go
$ ./main
Setting environment variable hello to world%
```

The output confirms successful code injection as the message appears when `os.Setenv()` is called.

There are many advanced features that is not covered in this example, you can refer to the existing rules in the `pkg/rules` directory for more information.
