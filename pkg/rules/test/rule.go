package test

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {

	api.NewRule("fmt", "Printf", "", "OnEnterPrintf1", "OnExitPrintf1").
		WithExecOrder(api.ExecOrderOutermost).
		WithFileDeps("long/sub/p2.go").
		WithPackageDep("google.golang.org/protobuf", "v1.34.0").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt",
		"Printf", "", "", `if internalFn!=nil { internalFn() } else { internalFn() }`).
		WithUseRaw(true).
		WithExecOrder(api.ExecOrderInnermost).
		WithFileDeps("long/sub/p1.go").
		WithRuleName("testrule").
		Register()

	api.NewFileRule("fmt", "long/sub/pp.go").
		WithRuleName("testrule").
		Register()

	api.NewStructRule("fmt", "pp", "myfield", "int").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "Fprintf", "", "", "n=7632; println(n);").WithUseRaw(true).
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "Unwrap", "", "onEnterUnwrap", "onExitUnwrap").
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "Do", "*Client", "onEnterClientDo", "onExitClientDo").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "Printf", "", "OnEnterPrintf2", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "NewRequest", "", "onEnterNewRequest", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "NewRequestWithContext", "", "onEnterNewRequestWithContext", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "internalFn", "", `println("GCMG")`, "").
		WithUseRaw(true).
		WithRuleName("testrule").
		Register()

	api.NewStructRule("fmt", "MyPoint", "y", "int").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "internalFn", "", `println(MyPoint{x:1024,y:512}.y)`, "").
		WithUseRaw(true).
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "NewRequest", "", "onEnterNewRequest1", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "Error", "*MaxBytesError", "onEnterMaxBytesError", "onExitMaxBytesError").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "newPrinter", "", "", "retVal0.myfield=0x7632; println(retVal0.myfield)").WithUseRaw(true).
		WithRuleName("testrule").
		Register()

	api.NewStructRule("net/http", "Request", "Should", "NotExist").
		WithRuleName("testrule").
		Register()

	api.NewFileRule("net/http", "long/sub/p3.go").
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "NewRequest", "", "", "onExitNewRequest").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "TestSkip", "", "onEnterTestSkip", "").
		WithRuleName("testrule").
		Register()

	api.NewFileRule("errors", "long/sub/p4.go").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "TestSkip", "", "", "onExitTestSkipOnly").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "TestSkip", "", "onEnterTestSkipOnly", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "Sprintf", "", "onEnterSprintf1", "onExitSprintf1").
		WithExecOrder(api.ExecOrderInner).
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "Sprintf", "", "onEnterSprintf2", "onExitSprintf2").
		WithExecOrder(api.ExecOrderOutermost).
		WithRuleName("testrule").
		Register()

	api.NewRule("fmt", "Sprintf", "", "onEnterSprintf3", "onExitSprintf3").
		WithExecOrder(api.ExecOrderInnermost).
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "p1", "", "onEnterP11", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "p1", "", "onEnterP12", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "p2", "", "", "onExitP21").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "p2", "", "", "onExitP22").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "p3", "", "onEnterP31", "onExitP31").
		WithRuleName("testrule").
		Register()

	// Test for version match
	api.NewRule("golang.org/x/time/rate", "Every", "",
		`println("BAD")`, "").
		WithUseRaw(true).
		WithVersion("[0.6.0,1.0.0)").
		WithRuleName("testrule").
		Register()

	api.NewRule("golang.org/x/time/rate", "Every", "",
		`println("BAD")`, "").
		WithUseRaw(true).
		WithVersion("[0.0.0,0.5.0)").
		WithRuleName("testrule").
		Register()

	api.NewRule("golang.org/x/time/rate", "Every", "",
		`println("GOOD")`, "").
		WithUseRaw(true).
		WithVersion("[0.0.0,999.999.999)").
		WithRuleName("testrule").
		Register()

	api.NewRule("golang.org/x/time/rate", "Every", "",
		`println("GOOD")`, "").
		WithUseRaw(true).
		WithVersion("[0.0.0,0.6.0)").
		WithRuleName("testrule").
		Register()

	api.NewRule("golang.org/x/time/rate", "Every", "",
		`println("BAD")`, "").
		WithUseRaw(true).
		WithVersion("[0.5.1, 0.6.0)").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "TestSkip2", "", "onEnterTestSkip2", "onExitTestSkip2").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "TestGetSet", "", "onEnterTestGetSet", "onExitTestGetSet").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "TestGetSetRecv", "*Recv", "onEnterTestGetSetRecv", "onExitTestGetSetRecv").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "OnlyRet", "", "", "onExitOnlyRet").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "OnlyArgs", "", "onEnterOnlyArgs", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("errors", "NilArg", "", "onEnterNilArg", "").
		WithRuleName("testrule").
		Register()
}
