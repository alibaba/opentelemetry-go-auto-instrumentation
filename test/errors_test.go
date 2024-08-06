package test

import (
	"regexp"
	"testing"
)

const ErrorsAppName = "errors-test"

func TestRunErrors(t *testing.T) {
	UseApp(ErrorsAppName)

	RunInstrument(t, "-debuglog", "-disablerules=fmt")
	stdout, stderr := RunApp(t, ErrorsAppName)
	ExpectContains(t, stdout, "wow")
	ExpectContains(t, stdout, "old:wow")
	ExpectContains(t, stdout, "ptr<nil>")
	ExpectNotContains(t, stdout, "val1024")
	ExpectContains(t, stdout, "val1298") // 0x512
	ExpectContains(t, stdout, "7632")
	ExpectContains(t, stdout, "4008208820")
	ExpectContains(t, stdout, "118888")
	ExpectContains(t, stdout, "0.001")
	ExpectContains(t, stderr, "2024 shanghai")
	ExpectContains(t, stdout, "2033 hangzhou")

	text := ReadInstrumentLog(t, "debug_fn_otel_inst_file_p4.go")
	re := regexp.MustCompile(".*OtelOnEnterTrampoline_TestSkip.*")
	matches := re.FindAllString(text, -1)
	if len(matches) < 1 {
		t.Fatalf("expecting at least one match")
	}
	re = regexp.MustCompile(".*OtelOnEnterTrampoline_p1.*")
	matches = re.FindAllString(text, -1)
	if len(matches) != 3 {
		t.Fatalf("expecting 3 matches")
	}
	re = regexp.MustCompile(".*OtelOnExitTrampoline_p2.*")
	matches = re.FindAllString(text, -1)
	if len(matches) != 3 {
		t.Fatalf("expecting 3 matches")
	}
}
