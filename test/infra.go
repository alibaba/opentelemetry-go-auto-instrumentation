package test

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const execName = "otel-go-auto-instrumentation"

func RunCmd(args []string) *exec.Cmd {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()
	stdoutFile := filepath.Join("stdout.log")
	stdout, _ := os.Create(stdoutFile)
	stderrFile := filepath.Join("stderr.log")
	stderr, _ := os.Create(stderrFile)

	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

func ReadInstrumentLog(t *testing.T, fileName string) string {
	path := filepath.Join(shared.TempBuildDir, shared.TInstrument, fileName)
	content, err := util.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func ReadPreprocessLog(t *testing.T, fileName string) string {
	path := filepath.Join(shared.TempBuildDir, shared.TPreprocess, fileName)
	content, err := util.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func readStdoutLog(t *testing.T) string {
	return readLog(t, "stdout.log")
}

func readStderrLog(t *testing.T) string {
	return readLog(t, "stderr.log")
}

func RunGoBuild(t *testing.T, args ...string) {
	util.Assert(pwd != "", "pwd is empty")
	cmd := RunCmd(append([]string{"go", "build"}, args...))
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func RunInstrumentFallible(t *testing.T, args ...string) {
	util.Assert(pwd != "", "pwd is empty")
	path := filepath.Join(filepath.Dir(pwd), execName)
	cmd := RunCmd(append([]string{path}, args...))
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected failure")
	}
}

func RunInstrument(t *testing.T, args ...string) {
	util.Assert(pwd != "", "pwd is empty")
	path := filepath.Join(filepath.Dir(pwd), execName)
	cmd := RunCmd(append([]string{path}, args...))
	err := cmd.Run()
	if err != nil {
		stderr := readStderrLog(t)
		log1 := ReadPreprocessLog(t, shared.DebugLogFile)
		log2 := ReadInstrumentLog(t, shared.DebugLogFile)
		text := fmt.Sprintf("failed to run instrument: %v\n", err)
		text += fmt.Sprintf("stderr: %v\n", stderr)
		text += fmt.Sprintf("preprocess: %v\n", log1)
		text += fmt.Sprintf("instrument: %v\n", log2)
		t.Fatal(text)
	}
}

func TBuildAppNoop(t *testing.T, appName string) {
	UseApp(appName)
	RunInstrument(t)
}

func BuildApp(t *testing.B, appName string) {
	UseApp(appName)
	RunInstrument(&testing.T{})
}

func BuildAppNoop(t *testing.B, appName string) {
	UseApp(appName)
	RunGoBuild(&testing.T{}, "-a" /*force rebuilding*/)
}

var pwd string

func UseApp(appName string) {
	if pwd == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current dir due to error", err.Error())
		}
		pwd = dir

	}
	err := os.Chdir(filepath.Join(pwd, appName))
	if err != nil {
		fmt.Println("Failed to chdir due to error", err.Error())
	}
}

func RunApp(t *testing.T, appName string, env ...string) (string, string) {
	cmd := RunCmd([]string{"./" + appName})
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)
	cmd.Env = append(cmd.Env, verifier.IS_IN_TEST+"=true")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	stdoutText := readStdoutLog(t)
	stderrText := readStderrLog(t)
	return stdoutText, stderrText
}

func FetchVersion(t *testing.T, dependency, version string) string {
	output, err := exec.Command("go", "get", "-u", dependency+"@"+version).Output()
	if err != nil {
		t.Fatal(output, err)
	}
	return string(output)
}

func readLog(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func ExpectStdoutContains(t *testing.T, expect string) {
	content := readStdoutLog(t)
	ExpectContains(t, content, expect)
}

func ExpectStderrContains(t *testing.T, expect string) {
	content := readStderrLog(t)
	ExpectContains(t, content, expect)
}

func ExpectInstrumentContains(t *testing.T, log string, rule string) {
	path := filepath.Join(shared.TempBuildDir, shared.TInstrument, log)
	content := readLog(t, path)
	ExpectContains(t, content, rule)
}

func ExpectInstrumentNotContains(t *testing.T, log string, rule string) {
	path := filepath.Join(shared.TempBuildDir, shared.TInstrument, log)
	content := readLog(t, path)
	ExpectNotContains(t, content, rule)
}

func ExpectPreprocessContains(t *testing.T, log string, rule string) {
	path := filepath.Join(shared.TempBuildDir, shared.TPreprocess, log)
	content := readLog(t, path)
	ExpectContains(t, content, rule)
}

func ExpectPreprocessNotContains(t *testing.T, log string, rule string) {
	path := filepath.Join(shared.TempBuildDir, shared.TPreprocess, log)
	content := readLog(t, path)
	ExpectNotContains(t, content, rule)
}

func ExpectContains(t *testing.T, text, expect string) {
	if !strings.Contains(text, expect) {
		t.Fatalf("text: %s, expect: %s", text, expect)
	}
}

func ExpectNotContains(t *testing.T, text, expect string) {
	if strings.Contains(text, expect) {
		t.Fatalf("text: %s, expect not: %s", text, expect)
	}
}

func ExpectSame(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func ExpectWhen(t *testing.T, prediction func() (res bool, msg string)) {
	if r, m := prediction(); !r {
		t.Fatalf(m)
	}
}

func ExpectContainsAllItem(t *testing.T, actualItems []string, expectedItems ...string) {
	expectedSet := make(map[string]*interface{})
	for _, item := range expectedItems {
		expectedSet[item] = nil
	}
	for _, item := range actualItems {
		delete(expectedSet, item)
	}
	if len(expectedSet) > 0 {
		sort.Strings(expectedItems)
		sort.Strings(actualItems)
		t.Fatalf("-- expected: [%s]\n-- actual: [%s]", strings.Join(expectedItems, ", "),
			strings.Join(actualItems, ", "))
	}
}

func ExpectContainsNothing(t *testing.T, actualItems []string) {
	if len(actualItems) > 0 {
		t.Fatalf("-- expected: []\n-- actual: [%s]", strings.Join(actualItems, ", "))
	}
}

func ExecMuzzle(t *testing.T, dependencyName, moduleName string, minVersion, maxVersion *Version) {
	versions, err := GetRandomVersion(3, dependencyName, minVersion, maxVersion)
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := os.ReadDir(moduleName)
	if err != nil {
		t.Fatal(err)
	}
	testVersions := make([]*Version, 0)
	for _, dir := range dirs {
		v, err := NewVersion(dir.Name())
		if err != nil {
			continue
		}
		testVersions = append(testVersions, v)
	}
	sort.Slice(testVersions, func(i, j int) bool {
		return testVersions[i].GreaterThan(testVersions[j])
	})
	for _, version := range versions {
		for _, testVersion := range testVersions {
			if version.GreaterThanOrEqual(testVersion) {
				t.Logf("testing on version %v\n", version.original)
				UseApp(moduleName + "/" + testVersion.original)
				FetchVersion(t, dependencyName, version.original)
				TBuildAppNoop(t, moduleName+"/"+testVersion.original)
				break
			}
		}
	}
}

func ExecLatestTest(t *testing.T, dependencyName, moduleName string, minVersion, maxVersion *Version, testFunc func(*testing.T, *Version)) {
	latestVersion, err := GetLatestVersion(dependencyName, minVersion, maxVersion)
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := os.ReadDir(moduleName)
	if err != nil {
		t.Fatal(err)
	}
	testVersions := make([]*Version, 0)
	for _, dir := range dirs {
		v, err := NewVersion(dir.Name())
		if err != nil {
			continue
		}
		testVersions = append(testVersions, v)
	}
	sort.Slice(testVersions, func(i, j int) bool {
		return testVersions[i].LessThan(testVersions[j])
	})
	latestTestVersion := testVersions[len(testVersions)-1]
	UseApp(moduleName + "/" + latestTestVersion.original)
	FetchVersion(t, dependencyName, latestVersion.original)
	testFunc(t, latestTestVersion)
}
