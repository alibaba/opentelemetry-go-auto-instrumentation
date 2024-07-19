# How it works

The automatic instrumentation solution is implemented around the concept of `Rule`.

Three types of rules defined in [ruledef.go](../api/ruledef.go) are supported now:

- `InstFuncRule`: describes how to instrument a specific go function.
- `InstStructRule`: describes how to instrument a specific go struct.
- `InstFileRule`: describes how to instrument a specific go file.

There are no mandatory restrictions on the scope of target code, which can be the project's own code, code from
dependent libraries, or even the go runtime.

The workflow could roughly be divided into two main phases:

- `Preprocess`: Analyze dependencies and select rules that should be used later.
- `Instrument`: Generate code based on rules and inject new code into source code.

TODO
