package test

import "testing"

func TestBuildProject(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)
	RunInstrument(t, "--", "-o", "default", "cmd/foo.go")
	RunInstrument(t, "--", "cmd/foo.go")
	RunInstrument(t, "--", "cmd/foo.go", "cmd/bar.go")
	RunInstrument(t, "--", "cmd")
}

func TestBuildProject2(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "--", ".")
	RunInstrument(t, "--", "")
}

func TestBuildProject3(t *testing.T) {
	const AppName = "build"
	UseApp(AppName)

	RunInstrument(t, "--", "m1")
	RunInstrumentFallible(t, "--", "m2") // not used in go.work
}
