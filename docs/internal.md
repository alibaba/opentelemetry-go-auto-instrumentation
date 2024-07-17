
# How does it work

The whole instrumentation framework is built on top of the concept `Rule`, which
describes how to instrument a specific method. By providing `Version`, `ImportPath`
`Function` and `ReceiverType`, we know exactly which method we are going to instrument.
Different kinds of rule could be defined to instrument in different ways. For
example, `InstFuncRule` is used to instrument a function, `InstStruct` is used to
instrument a struct, different rules could be combined to instrument the same
method in order to achieve the desired effect. If many rules are dedicated to
instrument the same method, the order of the rules is important, the `ExecOrder`
allows us to specify the order of the rules, the `ExecOrderOutermost` comes first,
the `ExecOrderInnermost` comes last.

There are some gory details in the implementation, it could roughly be divided 
into two parts: the *Preprocess* and *Instrument* phase. Given the user project,
it first determines available rules for the project, and the *Preprocess* phase 
is responsible for collecting all the dependencies and the available rules.
Once all preparations are done, we start the building by using `go build toolexec`
flag to hook the compilation action and the *Instrument* phase is responsible for
generating so-called **TJump** instrumentation code and injecting it into the source code.

# Trampoline Jump
**TJump** is short for `Trampoline Jump`, it also named `Trampoline Jump If` in some contexts.
We distinguish between three types of functions: RawFunc, TrampolineFunc, and
HookFunc. RawFunc is the original function that needs to be instrumented.
TrampolineFunc is the function that is generated to call the onEnter and
onExit hooks, it serves as a trampoline to the original function. HookFunc is
the function that is called at entrypoint and exitpoint of the RawFunc. The
so-called "Trampoline Jump" snippet is inserted at start of raw func, it is
guaranteed to be generated within one line to avoid confusing debugging, as
its name suggests, it jumps to the trampoline function from raw function.

There is a simple example to illustrate the concept, given the following RawFunc:
```go
func (c *Client) Do() {
    c.do()
}
```
We generate the following TJump, it jumps to the TrampolineFunc:

```go
func (c *Client) Do() {
    if ctx,skip := otel_trampoline_onenter(); skip {
        otel_trampoline_onexit(ctx)
        return
    } else {
        defer otel_trampoline_onexit(ctx)
    }
    c.do()
}
```

 Since trampoline-jump-if and trampoline functions are performance-critical,
 we are trying to optimize them as much as possible.

 The obvious optimization opportunities are cases when onEnter or onExit hooks
 are not present. For the latter case, we can replace the defer statement to
 empty statement, you might argue that we can remove the whole else block, but
 there might be more than one trampoline-jump-if in the same function, they are
 nested in the else block, i.e.

```go
  if ctx, skip := otel_trampoline_onenter(&arg); skip {
      otel_trampoline_onexit(ctx, &retval)
      return ...
  } else {
      ;
      ...
  }
```

 For the former case, it's a bit more complicated. We need to manually construct
 CallContext on the fly and pass it to onExit trampoline defer call and rewrite
 the whole condition to always false. The corresponding code snippet is

```go
  if false {
      ;
  } else {
      defer otel_trampoline_onexit(&CallContext{...}, &retval)
      ...
  }
```

 The if skeleton should be kept as is, otherwise inlining of trampoline-jump-if
 will not work. During compiling, the dce and sccp passes will remove the whole
 then block. That's not the whole story. We can further optimize the tjump iff
 the onEnter hook does not use SkipCall. In this case, we can rewrite condition
 of trampoline-jump-if to always false, remove return statement in then block,
 they are memory-aware and may generate memory SSA values during compilation.

```go
  if ctx,_ := otel_trampoline_onenter(&arg); false {
      ;
  } else {
      defer otel_trampoline_onexit(ctx, &retval)
      ...
  }
```

 The compiler responsible for hoisting the initialization statement out of the
 if skeleton, and the dce and sccp passes will remove the whole then block. All
 these trampoline functions looks as if they are executed sequentially.