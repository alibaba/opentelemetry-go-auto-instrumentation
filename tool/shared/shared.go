// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

const (
	GoBuildIgnoreComment = "//go:build ignore"
	GoModFile            = "go.mod"
	GoSumFile            = "go.sum"
	GoWorkSumFile        = "go.work.sum"
	DebugLogFile         = "debug.log"
	TInstrument          = "instrument"
	TPreprocess          = "preprocess"
)

func AssertGoBuild(args []string) {
	if len(args) < 2 {
		util.Assert(false, "empty go build command")
	}
	if !strings.Contains(args[0], "go") {
		util.Assert(false, "invalid go build command %v", args)
	}
	if args[1] != "build" {
		util.Assert(false, "invalid go build command %v", args)
	}
}

func IsCompileCommand(line string) bool {
	check := []string{"-o", "-p", "-buildid"}
	if util.IsWindows() {
		check = append(check, "compile.exe")
	} else if util.IsUnix() {
		check = append(check, "compile")
	} else {
		util.ShouldNotReachHere()
	}
	for _, id := range check {
		if !strings.Contains(line, id) {
			return false
		}
	}

	// PGO compile command is different from normal compile command, we
	// should skip it, otherwise the same package will be compiled twice
	// (one for PGO and one for normal), which finally leads to the same
	// rule being applied twice.
	if strings.Contains(line, "-pgoprofile=") {
		return false
	}
	return true
}

func GetTempBuildDir() string {
	if GetConf().InToolexec {
		return filepath.Join(TempBuildDir, TInstrument)
	} else {
		return filepath.Join(TempBuildDir, TPreprocess)
	}
}

func GetLogPath(name string) string {
	return filepath.Join(GetTempBuildDir(), name)
}

func GetInstrumentLogPath(name string) string {
	return filepath.Join(TempBuildDir, TInstrument, name)
}

func GetPreprocessLogPath(name string) string {
	return filepath.Join(TempBuildDir, TPreprocess, name)
}

func GetVarNameOfFunc(fn string) string {
	const varDeclSuffix = "Impl"
	fn = strings.Title(fn)
	return fn + varDeclSuffix
}

var packageRegexp = regexp.MustCompile(`(?m)^package\s+\w+`)

func RenamePackage(source, newPkgName string) string {
	source = packageRegexp.ReplaceAllString(source,
		fmt.Sprintf("package %s\n", newPkgName))
	return source
}

func RemoveGoBuildComment(text string) string {
	text = strings.ReplaceAll(text, GoBuildIgnoreComment, "")
	return text
}

func HasGoBuildComment(text string) bool {
	return strings.Contains(text, GoBuildIgnoreComment)
}

// GetGoModPath returns the absolute path of go.mod file, if any.
func GetGoModPath() (string, error) {
	// @@ As said in the comment https://github.com/golang/go/issues/26500, the
	// expected way to get go.mod should be go list -m -f {{.GoMod}}, but it does
	// not work well when go.work presents, we use go env GOMOD instead.
	//
	// go env GOMOD
	// The absolute path to the go.mod of the main module.
	// If module-aware mode is enabled, but there is no go.mod, GOMOD will be
	// os.DevNull ("/dev/null" on Unix-like systems, "NUL" on Windows).
	// If module-aware mode is disabled, GOMOD will be the empty string.
	cmd := exec.Command("go", "env", "GOMOD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get go.mod directory: %w\n%v",
			err, string(out))
	}
	path := strings.TrimSpace(string(out))
	return path, nil
}

// GetGoModDir returns the directory of go.mod file.
func GetGoModDir() (string, error) {
	gomod, err := GetGoModPath()
	if err != nil {
		return "", fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	projectDir := filepath.Dir(gomod)
	return projectDir, nil
}

// GetProjRootDir returns the root directory of the project. It's an alias of
// GetGoModDir in the current implementation.
func GetProjRootDir() (string, error) {
	return GetGoModDir()
}

// IsModPath checks if the provided module path is valid.
func IsModPath(path string) bool {
	if strings.Contains(path, "@") {
		pathOnly := strings.Split(path, "@")[0]
		return module.CheckPath(pathOnly) == nil
	}
	return module.CheckPath(path) == nil
}

func IsGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

func IsGoModFile(path string) bool {
	return strings.HasSuffix(path, GoModFile)
}

func IsGoSumFile(path string) bool {
	return strings.HasSuffix(path, "go.sum")
}

func IsGoTestFile(path string) bool {
	return strings.HasSuffix(path, "_test.go")
}

func IsExistGoMod() (bool, error) {
	gomod, err := GetGoModPath()
	if err != nil {
		return false, fmt.Errorf("failed to get go.mod path: %w", err)
	}
	if gomod == "" {
		return false, errors.New("failed to get go.mod path: not module-aware")
	}
	return strings.HasSuffix(gomod, GoModFile), nil
}

func HashStruct(st interface{}) (uint64, error) {
	bs, err := json.Marshal(st)
	if err != nil {
		return 0, err
	}
	hasher := fnv.New64a()
	_, err = hasher.Write(bs)
	if err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

func MakePublic(name string) string {
	return strings.Title(name)
}

func InPreprocess() bool {
	return !GetConf().InToolexec
}

func InInstrument() bool {
	return GetConf().InToolexec
}

func GuaranteeInPreprocess() {
	util.Assert(!GetConf().InToolexec, "not in preprocess stage")
}

func GuaranteeInInstrument() {
	util.Assert(GetConf().InToolexec, "not in instrument stage")
}

// splitVersionRange splits the version range into two parts, start and end.
func splitVersionRange(vr string) (string, string) {
	util.Assert(strings.Contains(vr, ","), "invalid version range format")
	util.Assert(strings.Contains(vr, "["), "invalid version range format")
	util.Assert(strings.Contains(vr, ")"), "invalid version range format")

	start := vr[1:strings.Index(vr, ",")]
	end := vr[strings.Index(vr, ",")+1 : len(vr)-1]
	return "v" + start, "v" + end
}

var versionRegexp = regexp.MustCompile(`@v\d+\.\d+\.\d+(-.*?)?/`)

func ExtractVersion(path string) string {
	// Unify the path to Unix style
	path = filepath.ToSlash(path)
	version := versionRegexp.FindString(path)
	if version == "" {
		return ""
	}
	// Extract version number from the string
	return version[1 : len(version)-1]
}

// MatchVersion checks if the version string matches the version range in the
// rule. The version range is in format [start, end), where start is inclusive
// and end is exclusive. If the rule version string is empty, it always matches.
func MatchVersion(version string, ruleVersion string) (bool, error) {
	// Fast path, always match if the rule version is not specified
	if ruleVersion == "" {
		return true, nil
	}
	// Check if both rule version and package version are in sane
	if !strings.Contains(version, "v") {
		return false, fmt.Errorf("invalid version %v %v",
			version, ruleVersion)
	}
	if !strings.Contains(ruleVersion, "[") ||
		!strings.Contains(ruleVersion, ")") ||
		!strings.Contains(ruleVersion, ",") ||
		strings.Contains(ruleVersion, "v") {
		return false, fmt.Errorf("invalid version format in rule %v",
			ruleVersion)
	}
	// Remove extra whitespace from the rule version string
	ruleVersion = strings.ReplaceAll(ruleVersion, " ", "")

	// Compare the version with the rule version, the rule version is in the
	// format [start, end), where start is inclusive and end is exclusive
	ruleVersionStart, ruleVersionEnd := splitVersionRange(ruleVersion)
	if semver.Compare(version, ruleVersionStart) >= 0 &&
		semver.Compare(version, ruleVersionEnd) < 0 {
		return true, nil
	}
	return false, nil
}

// SplitCmds splits the command line by space, but keep the quoted part as a
// whole. For example, "a b" c will be split into ["a b", "c"].
func SplitCmds(input string) []string {
	var args []string
	var inQuotes bool
	var arg strings.Builder

	for i := 0; i < len(input); i++ {
		c := input[i]

		if c == '"' {
			inQuotes = !inQuotes
			continue
		}

		if c == ' ' && !inQuotes {
			if arg.Len() > 0 {
				args = append(args, arg.String())
				arg.Reset()
			}
			continue
		}

		arg.WriteByte(c)
	}

	if arg.Len() > 0 {
		args = append(args, arg.String())
	}

	// Fix the escaped backslashes on Windows
	if util.IsWindows() {
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, `\\`, `\`)
		}
	}
	return args
}
