package internal

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/golang-collections/collections/set"
	"os"
	"os/exec"
	"otel-auto-instrumentation/internal/shared"
	"path/filepath"
	"strings"
)

const workingDirName = ".otel-auto-instrumentation"
const workingDirEnv = "OTEL_AUTO_INSTRUMENTATION_WORKING_DIRECTORY"
const logFileName = "log"
const dryRunLogName = "dry-run-log"

var dryRunLogPath string
var dryRunLog *os.File

//go:embed setup.go
var setupSnippet string

func Run() (err error) {
	err = initContext()
	if err != nil {
		return fmt.Errorf("failed to init context: %w", err)
	}

	if !shared.InToolexec {
		err = dryRun()
		if err != nil {
			return fmt.Errorf("failed to dry run: %w", err)
		}

		pkgs, err := readPackages()
		if err != nil {
			return fmt.Errorf("failed to read packages: %w", err)
		}

		err = injectImports()
		if err != nil {
			return fmt.Errorf("failed to inject imports: %w", err)
		}

		err = injectSnippets(pkgs)
		if err != nil {
			return fmt.Errorf("failed to inject snippets: %w", err)
		}

		err = runGet()
		if err != nil {
			return fmt.Errorf("failed to run go get: %w", err)
		}

		err = runBuild()
		if err != nil {
			return fmt.Errorf("failed to run go build: %w", err)
		}
		return nil
	} else {
		args := os.Args[2:]
		if strings.HasSuffix(os.Args[2], "compile") {
			pkg, od := readPkgAndOutputDir(os.Args)
			r, exist := rules[pkg]
			if exist {
				args, err = apply(r, args, od)
				if err != nil {
					return fmt.Errorf("failed to apply rule '%s'", r.Name)
				}
			}
		}
		return runCmd(args)
	}
}

func initContext() (err error) {
	if shared.InToolexec {
		wd := os.Getenv(workingDirEnv)
		logPath := filepath.Join(wd, logFileName)
		if len(wd) == 0 {
			return fmt.Errorf("cannot find working direcotory")
		}
		shared.WorkingDir = wd
		shared.LogPath = logPath
	} else {
		s, err := os.Stat(workingDirName)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat working directory: %w", err)
		}

		if s != nil {
			err = os.RemoveAll(workingDirName)
			if err != nil {
				return fmt.Errorf("failed to remove existing working directory: %w", err)
			}
		}

		err = os.Mkdir(workingDirName, 0755)
		if err != nil {
			return fmt.Errorf("failed to make working directory: %w", err)
		}

		wd, err := filepath.Abs(workingDirName)
		if err != nil {
			return err
		}

		logPath := filepath.Join(wd, logFileName)
		logFile, err := os.Create(logPath)
		if err != nil {
			return fmt.Errorf("failed to create log: %w", err)
		}
		defer logFile.Close()

		dryRunLogPath = filepath.Join(wd, dryRunLogName)
		if dryRunLog, err = os.Create(dryRunLogPath); err != nil {
			return fmt.Errorf("failed to create log: %w", err)
		}

		if err = os.Setenv(workingDirEnv, wd); err != nil {
			return err
		}

		shared.WorkingDir = wd
		shared.LogPath = logPath
	}
	return nil
}

func dryRun() error {
	cmd := exec.Command("go", "build", "-a", "-x", "-n")
	cmd.Stdin = os.Stdin
	cmd.Stdout = dryRunLog
	cmd.Stderr = dryRunLog
	return cmd.Run()
}

func readPackages() (pkgs *set.Set, err error) {
	pkgs = set.New()

	// reopen
	if dryRunLog, err = os.Open(dryRunLogPath); err != nil {
		return nil, err
	}
	defer dryRunLog.Close()

	scanner := bufio.NewScanner(dryRunLog)

	pf := "packagefile"
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, pf) {
			i := strings.Index(line, "=")
			if i != -1 {
				pkgs.Insert(line[len(pf)+1 : i])
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return pkgs, nil
}

func injectImports() (err error) {
	imports := []string{
		"go.opentelemetry.io/otel",
	}

	c := "package main\n"
	for _, i := range imports {
		c += "import _ \"" + i + "\"\n"
	}

	f, err := os.OpenFile("otel_imports.go", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(c)
	return err
}

func injectSnippets(pkgs *set.Set) (err error) {
	pkgs.Do(func(e interface{}) {
		if err != nil {
			return
		}
		pkg := e.(string)
		rule, exist := rules[pkg]
		if exist {
			f, _ := os.OpenFile("otel_"+rule.Name+"_snippets.go", os.O_CREATE|os.O_WRONLY, 0644)
			_, err = f.WriteString(rule.SetupSnippet)
		}
	})

	if err != nil {
		return err
	}

	snippet := strings.Replace(setupSnippet, "package internal", "package main", 1)
	f, _ := os.OpenFile("otel_setup_snippets.go", os.O_CREATE|os.O_WRONLY, 0644)
	_, err = f.WriteString(snippet)
	return err
}

func runGet() error {
	cmd := exec.Command("go", "get")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runBuild() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	toolexecArg := "-toolexec=" + exe + " -in-toolexec"
	cmd := exec.Command("go", "build", toolexecArg, "-a", "-work", "-x")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func readPkgAndOutputDir(cmd []string) (string, string) {
	var pkg, outputDir string
	for i, v := range cmd {
		if v == "-p" {
			pkg = cmd[i+1]
		} else if v == "-o" {
			outputDir = filepath.Dir(cmd[i+1])
		}
	}
	return pkg, outputDir
}

func runCmd(args []string) error {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
