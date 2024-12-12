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
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type RunPhase string

const (
	PInvalid    = "invalid"
	PPreprocess = "preprocess"
	PInstrument = "instrument"
	PConfigure  = "configure"
)

var rp RunPhase = "bad"

func SetRunPhase(phase RunPhase) {
	rp = phase
}

func GetRunPhase() RunPhase {
	return rp
}

func (rp RunPhase) String() string {
	return string(rp)
}

func InPreprocess() bool {
	return rp == PPreprocess
}

func InInstrument() bool {
	return rp == PInstrument
}

func InConfigure() bool {
	return rp == PConfigure
}

func GuaranteeInPreprocess() {
	Assert(rp == PPreprocess, "not in preprocess stage")
}

func GuaranteeInInstrument() {
	Assert(rp == PInstrument, "not in instrument stage")
}

func GuaranteeInConfigure() {
	Assert(rp == PConfigure, "not in configure stage")
}

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

var recordedRand = make(map[string]bool)

// RandomString generates a globally unique random string of length n
func RandomString(n int) string {
	for {
		var letters = []rune("0123456789")
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		s := string(b)
		// Random suffix collision? Reroll until we get a unique one
		if _, ok := recordedRand[s]; !ok {
			recordedRand[s] = true
			return s
		}
	}
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

func RunCmdOutput(args ...string) (string, error) {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			LogFatal("failed to close file %s: %v", sourceFile.Name(), err)
		}
	}(sourceFile)

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			LogFatal("failed to close file %s: %v", destFile.Name(), err)
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
			LogFatal("failed to close file %s: %v", file.Name(), err)
		}
	}(file)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

func WriteFile(filePath string, content string) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogFatal("failed to close file %s: %v", file.Name(), err)
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
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Dont list files under hidden directories
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}
	err := filepath.Walk(dir, walkFn)
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

func CopyDir(src string, dst string) error {
	// Get the properties of the source directory
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	// Read the contents of the source directory
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	// Iterate through each entry in the source directory
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
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

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin"
}

func PhaseTimer(name string) func() {
	start := time.Now()
	return func() {
		Log("%s took %f s", name, time.Since(start).Seconds())
	}
}

// StringDedup removes duplicate strings in a slice and returns a new slice
// in original order.
func StringDedup(s []string) []string {
	m := make(map[string]struct{})
	var r []string
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

func GetToolName() string {
	// Get the path of the current executable
	ex, err := os.Executable()
	if err != nil {
		LogFatal("failed to get executable: %v", err)
		os.Exit(0)
	}
	return filepath.Base(ex)
}
