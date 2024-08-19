package test

import "testing"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("gin-test", "gin", "", "", "1.21", "", TestBasicGin),
	)
}

func TestBasicGin(t *testing.T, env ...string) {
	UseApp("gin")
	RunInstrument(t, "-debuglog", "--", "gin.go")
	RunApp(t, "gin", env...)
}
