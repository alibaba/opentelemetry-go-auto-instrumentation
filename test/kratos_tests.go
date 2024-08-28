package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("kratos-test", "kratos", "", "", "1.21", "", TestBasicKratos),
	)
}

func TestBasicKratos(t *testing.T, env ...string) {
	UseApp("kratos")
	RunInstrument(t, "-debuglog")
	RunApp(t, "kratos", env...)
}
