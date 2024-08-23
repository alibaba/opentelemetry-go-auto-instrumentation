package test

import (
	"testing"
)

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("echo-test", "echo", "4.10.0", "4.12.1", "1.21", "", TestBasicEcho),
	)
}

func TestBasicEcho(t *testing.T, env ...string) {
	UseApp("echo")
	RunInstrument(t, "-debuglog", "--", "echo.go")
	RunApp(t, "echo", env...)
}
