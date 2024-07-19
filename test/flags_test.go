package test

import (
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

func TestFlags(t *testing.T) {
	const AppName = "flags"
	UseApp(AppName)

	RunInstrument(t, "-debuglog", "-disablerules=")
	ExpectInstrumentContains(t, util.DebugLogFile, "fmt@")

	RunInstrument(t, "-debuglog", "-disablerules=fmt,net/http")
	ExpectInstrumentNotContains(t, util.DebugLogFile, "fmt@")
	ExpectInstrumentNotContains(t, util.DebugLogFile, "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=*")
	ExpectInstrumentNotContains(t, util.DebugLogFile, "fmt@")
	ExpectInstrumentNotContains(t, util.DebugLogFile, "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=testrule")
	ExpectInstrumentNotContains(t, util.DebugLogFile, "fmt@")

	RunInstrumentFallible(t, "-debuglog", "--", "-thisisnotvalid")
	ExpectPreprocessContains(t, util.DebugLogFile, "failed to")

	RunInstrument(t, "-version")
	ExpectStdoutContains(t, "version")

	RunInstrumentFallible(t, "-debuglog", "--", "notevenaflag")
	ExpectPreprocessContains(t, util.DebugLogFile, "failed to")

	RunInstrument(t, "-debuglog", "-verbose",
		"--",
		`-ldflags=-X main.Placeholder=replaced`)
	_, stderr := RunApp(t, AppName)
	ExpectContains(t, stderr, "placeholder:replaced")
}
