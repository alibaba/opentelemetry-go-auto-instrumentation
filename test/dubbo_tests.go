package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("dubbo-test", "dubbo", "", "", "1.21", "", TestBasicDubbo),
	)
}

func TestBasicDubbo(t *testing.T, env ...string) {
	UseApp("dubbo")
	RunInstrument(t, "-debuglog")
	RunApp(t, "dubbo", env...)
}
