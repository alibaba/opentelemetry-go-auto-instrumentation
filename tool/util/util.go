package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"golang.org/x/mod/modfile"
)

const GoBuildIgnoreComment = "//go:build ignore"

const GoModFile = "go.mod"

const (
	TInstrument = "instrument"
	TPreprocess = "preprocess"
)

var Guarantee = Assert // More meaningful name:)

func Assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}

func Unimplemented() {
	panic("unimplemented")
}

func UnimplementedT(msg string) {
	panic("unimplemented: " + msg)
}

func ShouldNotReachHere() {
	panic("should not reach here")
}

func ShouldNotReachHereT(msg string) {
	panic("should not reach here: " + msg)
}

func IsCompileCommand(line string) bool {
	return strings.Contains(line, "compile -o") &&
		strings.Contains(line, "buildid")
}

func GetLogPath(name string) string {
	if shared.InToolexec {
		return filepath.Join(shared.TempBuildDir, TInstrument, name)
	} else {
		return filepath.Join(shared.TempBuildDir, TPreprocess, name)
	}
}

func GetInstrumentLogPath(name string) string {
	return filepath.Join(shared.TempBuildDir, TInstrument, name)
}

func GetPreprocessLogPath(name string) string {
	return filepath.Join(shared.TempBuildDir, TPreprocess, name)
}

func GetVarNameOfFunc(fn string) string {
	const varDeclSuffix = "Impl"
	fn = strings.Title(fn)
	return fn + varDeclSuffix
}

func SaveDebugFile(prefix string, path string) {
	targetName := filepath.Base(path)
	Assert(IsGoFile(targetName), "sanity check")
	counterpart := GetLogPath("debug_" + prefix + targetName)
	_ = CopyFile(path, counterpart)
}

var packageRegexp = regexp.MustCompile(`(?m)^package\s+\w+`)

func RenamePackage(source, newPkgName string) string {
	source = packageRegexp.ReplaceAllString(source, fmt.Sprintf("package %s\n", newPkgName))
	return source
}

func RemoveGoBuildComment(text string) string {
	text = strings.ReplaceAll(text, GoBuildIgnoreComment, "")
	return text
}

func HasGoBuildComment(text string) bool {
	return strings.Contains(text, GoBuildIgnoreComment)
}

// getGoModPath returns the absolute path of go.mod file, if any.
func getGoModPath() (string, error) {
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
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	path := strings.TrimSpace(string(out))
	return path, nil
}

func IsGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

func IsExistGoMod() (bool, error) {
	gomod, err := getGoModPath()
	if err != nil {
		return false, fmt.Errorf("failed to get go.mod path: %w", err)
	}
	if gomod == "" {
		return false, errors.New("failed to get go.mod path: not module-aware")
	}
	return strings.HasSuffix(gomod, GoModFile), nil
}

// getModuleName returns the module name of the project by parsing go.mod file.
func getModuleName(gomod string) (string, error) {
	data, err := ReadFile(gomod)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	modFile, err := modfile.Parse("go.mod", []byte(data), nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse go.mod: %w", err)
	}

	moduleName := modFile.Module.Mod.Path
	return moduleName, nil
}

func GetImportPathOf(dirName string) (string, error) {
	Assert(dirName != "", "dirName is empty")
	// Get absolute path of current working directory
	workingDir, err := filepath.Abs(".")
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	// Get absolute path of go.mod directory
	gomod, err := getGoModPath()
	if err != nil {
		return "", fmt.Errorf("failed to get go.mod directory: %w", err)
	}
	projectDir := filepath.Dir(gomod)
	// Replace go.mod directory with module name
	moduleName, err := getModuleName(gomod)
	if err != nil {
		return "", fmt.Errorf("failed to get module name: %w", err)
	}
	moduleName = strings.Replace(workingDir, projectDir, moduleName, 1)
	return moduleName + "/" + dirName, nil
}

func RandomString(n int) string {
	var letters = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func StringQuote(s string) string {
	return `"` + s + `"`
}

func RunCmd(args ...string) error {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			log.Fatalf("failed to close file %s: %v", sourceFile.Name(), err)
		}
	}(sourceFile)

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			log.Fatalf("failed to close file %s: %v", destFile.Name(), err)
		}
	}(destFile)

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file %s: %v", file.Name(), err)
		}
	}(file)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

func WriteStringToFile(filePath string, content string) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file %s: %v", file.Name(), err)
		}
	}(file)

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func ListFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ListFilesFlat(dir string) ([]string, error) {
	// no recursive
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}
	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths, nil
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

func IsDebugMode() bool {
	return shared.Debug
}

func IsProductMode() bool {
	return !shared.Debug
}

func InPreprocess() bool {
	return !shared.InToolexec
}

func InInstrument() bool {
	return shared.InToolexec
}

func GuaranteeInPreprocess() {
	Assert(!shared.InToolexec, "not in preprocess stage")
}

func GuaranteeInInstrument() {
	Assert(shared.InToolexec, "not in instrument stage")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
