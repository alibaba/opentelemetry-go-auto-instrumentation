package util

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
