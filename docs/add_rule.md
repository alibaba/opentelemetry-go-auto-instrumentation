# How to add a new rule for desired framework
To add a new rule for a desired framework, you need to follow the steps below:

1. Mkdir in **pkg/rules** directory and create a new **rule.go** file
2. Register a new rule that specifies the framework import path, function
    name, and the function signature you want to instrument, and onEnter/onExit
    functions you want to insert before/after calling original function, respectively, in the **rule.go** file, as shown below:
    ```go
    package rules

    func init() {
        api.NewRule("framework_name", "function_name", "function_signature", "onEnterFunc", "onExitFunc").Register()
    }
    ```
3. Implement the onEnter/onExit functions in the same directory with the **rule.go** file, as shown below:
```go
// myhook.go
package rules

import framework_name

func onEnterFunc(ctx *framework_name.Context, arg1 int, arg2 bool) {
    println("onEnter")
}

func onExitFunc(ctx *framework_name.Context, ret1 string) {
    println("onExit")
}
```
4. Import your new rule package by adding import statement within `rule_enabler.go`

There are some concrete examples in the **pkg/rules** directory, such as **pkg/rules/test**. You can refer to them for more details.