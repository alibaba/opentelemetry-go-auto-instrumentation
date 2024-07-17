package test

import (
	"testing"
)

func TestFlags(t *testing.T) {
	const AppName = "flags"
	UseApp(AppName)

	RunInstrument(t, "-debuglog", "-disablerules=")
	ExpectInstrumentContains(t, "debug.log", "fmt@")

	RunInstrument(t, "-debuglog", "-disablerules=fmt,net/http")
	ExpectInstrumentNotContains(t, "debug.log", "fmt@")
	ExpectInstrumentNotContains(t, "debug.log", "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=*")
	ExpectInstrumentNotContains(t, "debug.log", "fmt@")
	ExpectInstrumentNotContains(t, "debug.log", "net/http@")

	RunInstrument(t, "-debuglog", "-disablerules=testrule")
	ExpectInstrumentNotContains(t, "debug.log", "fmt@")

	RunInstrumentFallible(t, "-debuglog", "--", "-thisisnotvalid")
	ExpectPreprocessContains(t, "debug.log", "failed to")

	RunInstrument(t, "-version")
	ExpectStdoutContains(t, "version")

	RunInstrumentFallible(t, "-debuglog", "--", "notevenaflag")
	ExpectPreprocessContains(t, "debug.log", "failed to")

	RunInstrument(t, "-debuglog", "-verbose",
		"--",
		`-ldflags=-X main.Placeholder=replaced`)
	_, stderr := RunApp(t, AppName)
	ExpectContains(t, stderr, "placeholder:replaced")
}
