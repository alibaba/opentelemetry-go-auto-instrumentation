package test

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"testing"
)

func TestFlags(t *testing.T) {
	const AppName = "flags"
	UseApp(AppName)

	RunInstrument(t, "-debuglog", "-disablerules=")
	ExpectInstrumentContains(t, shared.DebugLogFile, "fmt@")

	RunInstrument(t, "-debuglog", "-disablerules=fmt,net/http")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "fmt@")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=*")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "fmt@")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=testrule")
	ExpectInstrumentNotContains(t, shared.DebugLogFile, "fmt@")

	RunInstrumentFallible(t, "-debuglog", "--", "-thisisnotvalid")
	ExpectPreprocessContains(t, shared.DebugLogFile, "failed to")

	RunInstrument(t, "-version")
	ExpectStdoutContains(t, "version")

	RunInstrumentFallible(t, "-debuglog", "--", "notevenaflag")
	ExpectPreprocessContains(t, shared.DebugLogFile, "failed to")

	RunInstrument(t, "-debuglog", "-verbose",
		"--",
		`-ldflags=-X main.Placeholder=replaced`)
	_, stderr := RunApp(t, AppName)
	ExpectContains(t, stderr, "placeholder:replaced")
}
