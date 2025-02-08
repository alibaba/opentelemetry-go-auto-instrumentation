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

package util

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

const (
	GoBuildIgnoreComment = "//go:build ignore"
	GoModFile            = "go.mod"
	GoSumFile            = "go.sum"
	GoWorkSumFile        = "go.work.sum"
	DebugLogFile         = "debug.log"
	TempBuildDir         = ".otel-build"
	BuildConfFile        = "build_conf.json"
)

const (
	BuildPattern    = "-p"
	BuildGoVer      = "-goversion"
	BuildPgoProfile = "-pgoprofile"
	BuildModeVendor = "-mod=vendor"
	BuildModeMod    = "-mod=mod"
	BuildWork       = "-work"
)

func AssertGoBuild(args []string) {
	if len(args) < 2 {
		Assert(false, "empty go build command")
	}
	if !strings.Contains(args[0], "go") {
		Assert(false, "invalid go build command %v", args)
	}
	if args[1] != "build" {
		Assert(false, "invalid go build command %v", args)
	}
}

func IsCompileCommand(line string) bool {
	check := []string{"-o", "-p", "-buildid"}
	if IsWindows() {
		check = append(check, "compile.exe")
	} else if IsUnix() {
		check = append(check, "compile")
	} else {
		ShouldNotReachHere()
	}

	// Check if the line contains all the required fields
	for _, id := range check {
		if !strings.Contains(line, id) {
			return false
		}
	}

	// @@PGO compile command is different from normal compile command, we
	// should skip it, otherwise the same package will be compiled twice
	// (one for PGO and one for normal), which finally leads to the same
	// rule being applied twice.
	if strings.Contains(line, BuildPgoProfile) {
		return false
	}
	return true
}

func GetTempBuildDir() string {
	return filepath.Join(TempBuildDir, GetRunPhase().String())
}

func GetTempBuildDirWith(name string) string {
	return filepath.Join(TempBuildDir, name)
}

func GetLogPath(name string) string {
	return filepath.Join(GetTempBuildDir(), name)
}

func GetInstrumentLogPath(name string) string {
	return filepath.Join(TempBuildDir, PInstrument, name)
}

func GetPreprocessLogPath(name string) string {
	return filepath.Join(TempBuildDir, PPreprocess, name)
}

func GetConfigureLogPath(name string) string {
	return filepath.Join(TempBuildDir, PConfigure, name)
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

func HashStruct(st interface{}) (uint64, error) {
	bs, err := json.Marshal(st)
	if err != nil {
		return 0, errc.New(errc.ErrInvalidJSON, err.Error())
	}
	hasher := fnv.New64a()
	_, err = hasher.Write(bs)
	if err != nil {
		return 0, errc.New(errc.ErrInternal, err.Error())
	}
	return hasher.Sum64(), nil
}

func MakePublic(name string) string {
	return strings.Title(name)
}

// splitVersionRange splits the version range into two parts, start and end.
func splitVersionRange(vr string) (string, string) {
	Assert(strings.Contains(vr, ","), "invalid version range format")
	Assert(strings.Contains(vr, "["), "invalid version range format")
	Assert(strings.Contains(vr, ")"), "invalid version range format")

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
		return false, errc.New(errc.ErrMatchRule,
			fmt.Sprintf("invalid version %v", version))
	}
	if !strings.Contains(ruleVersion, "[") ||
		!strings.Contains(ruleVersion, ")") ||
		!strings.Contains(ruleVersion, ",") ||
		strings.Contains(ruleVersion, "v") {
		return false, errc.New(errc.ErrMatchRule,
			fmt.Sprintf("invalid rule version %v", ruleVersion))
	}
	// Remove extra whitespace from the rule version string
	ruleVersion = strings.ReplaceAll(ruleVersion, " ", "")

	// Compare the version with the rule version, the rule version is in the
	// format [start, end), where start is inclusive and end is exclusive
	// and start or end can be omitted, which means the range is open-ended.
	ruleVersionStart, ruleVersionEnd := splitVersionRange(ruleVersion)
	switch {
	case ruleVersionStart != "v" && ruleVersionEnd != "v":
		// Full version range
		if semver.Compare(version, ruleVersionStart) >= 0 &&
			semver.Compare(version, ruleVersionEnd) < 0 {
			return true, nil
		}
	case ruleVersionStart == "v":
		// Only end is specified
		Assert(ruleVersionEnd != "v", "sanity check")
		if semver.Compare(version, ruleVersionEnd) < 0 {
			return true, nil
		}
	case ruleVersionEnd == "v":
		// Only start is specified
		Assert(ruleVersionStart != "v", "sanity check")
		if semver.Compare(version, ruleVersionStart) >= 0 {
			return true, nil
		}
	default:
		return false, errc.New(errc.ErrMatchRule,
			fmt.Sprintf("invalid rule version range %v", ruleVersion))
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
	if IsWindows() {
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, `\\`, `\`)
		}
	}
	return args
}
