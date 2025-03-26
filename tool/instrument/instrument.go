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

package instrument

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

// -----------------------------------------------------------------------------
// Instrument
//
// The instrument package is used to instrument the source code according to the
// predefined rules. It finds the rules that match the project dependencies and
// applies the rules to the dependencies one by one.

type RuleProcessor struct {
	packageName     string
	workDir         string
	target          *dst.File       // The target file to be instrumented
	parser          *util.AstParser // The parser for the target file
	compileArgs     []string
	rule2Suffix     map[*resource.InstFuncRule]string
	rawFunc         *dst.FuncDecl
	exact           bool // If the rule is exact match with target function
	onEnterHookFunc *dst.FuncDecl
	onExitHookFunc  *dst.FuncDecl
	varDecls        []dst.Decl
	relocated       map[string]string
	trampolineJumps []*TJump // Optimization candidates
	callCtxDecl     *dst.GenDecl
	callCtxMethods  []*dst.FuncDecl
}

func newRuleProcessor(args []string, pkgName string) *RuleProcessor {
	// Read compilation output directory
	var outputDir string
	for i, v := range args {
		if v == "-o" {
			outputDir = filepath.Dir(args[i+1])
			break
		}
	}
	util.Assert(outputDir != "", "sanity check")
	// Create a new rule processor
	rp := &RuleProcessor{
		packageName: pkgName,
		workDir:     outputDir,
		target:      nil,
		compileArgs: args,
		rule2Suffix: make(map[*resource.InstFuncRule]string),
		relocated:   make(map[string]string),
	}
	return rp
}

func (rp *RuleProcessor) addDecl(decl dst.Decl) {
	rp.target.Decls = append(rp.target.Decls, decl)
}

func (rp *RuleProcessor) removeDeclWhen(pred func(dst.Decl) bool) dst.Decl {
	for i, decl := range rp.target.Decls {
		if pred(decl) {
			rp.target.Decls = append(rp.target.Decls[:i], rp.target.Decls[i+1:]...)
			return decl
		}
	}
	return nil
}

func (rp *RuleProcessor) setRelocated(name, target string) {
	rp.relocated[name] = target
}

func (rp *RuleProcessor) tryRelocated(name string) string {
	if target, ok := rp.relocated[name]; ok {
		return target
	}
	return name
}

func (rp *RuleProcessor) addCompileArg(newArg string) {
	rp.compileArgs = append(rp.compileArgs, newArg)
}

func haveSameSuffix(s1, s2 string) bool {
	minLength := len(s1)
	if len(s2) < minLength {
		minLength = len(s2)
	}
	for i := 1; i <= minLength; i++ {
		if s1[len(s1)-i] != s2[len(s2)-i] {
			return false
		}
	}
	return true
}

func (rp *RuleProcessor) replaceCompileArg(newArg string, pred func(string) bool) error {
	variant := ""
	for i, arg := range rp.compileArgs {
		// Use absolute file path of the compile argument to compare with the
		// instrumented file(path), which is also an absolute path
		arg, err := filepath.Abs(arg)
		if err != nil {
			return errc.New(errc.ErrAbsPath, err.Error())
		}
		if pred(arg) {
			rp.compileArgs[i] = newArg
			// Relocate the replaced file to new target, any rules targeting the
			// replaced file should be updated to target the new file as well
			rp.setRelocated(arg, newArg)
			return nil
		}
		if haveSameSuffix(arg, newArg) {
			variant = arg
		}
	}
	if variant == "" {
		variant = fmt.Sprintf("%v", rp.compileArgs)
	}
	msg := fmt.Sprintf("expect %s, actual %s", newArg, variant)
	return errc.New(errc.ErrInstrument, msg)
}

func (rp *RuleProcessor) saveDebugFile(path string) {
	escape := func(s string) string {
		dirName := strings.ReplaceAll(s, "/", "_")
		dirName = strings.ReplaceAll(dirName, ".", "_")
		return dirName
	}
	dest := filepath.Base(path)
	util.Assert(rp.packageName != "", "sanity check")
	dest = filepath.Join(escape(rp.packageName), dest)
	dest = util.GetInstrumentLogPath(dest)
	err := os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	if err != nil { // error is tolerable here
		util.Log("failed to create debug file directory %s: %v", dest, err)
		return
	}
	err = util.CopyFile(path, dest)
	if err != nil { // error is tolerable here
		util.Log("failed to save debug file %s: %v", dest, err)
	}
}

func (rp *RuleProcessor) applyRules(bundle *resource.RuleBundle) (err error) {
	// Apply file instrument rules first
	err = rp.applyFileRules(bundle)
	if err != nil {
		err = errc.Adhere(err, "package", bundle.ImportPath)
		return err
	}

	err = rp.applyStructRules(bundle)
	if err != nil {
		err = errc.Adhere(err, "package", bundle.ImportPath)
		return err
	}

	err = rp.applyFuncRules(bundle)
	if err != nil {
		err = errc.Adhere(err, "package", bundle.ImportPath)
		return err
	}

	return nil
}

func matchImportPath(importPath string, args []string) bool {
	for _, arg := range args {
		if arg == importPath {
			return true
		}
	}
	return false
}

func compileRemix(bundle *resource.RuleBundle, args []string) error {
	rp := newRuleProcessor(args, bundle.PackageName)
	err := rp.applyRules(bundle)
	if err != nil {
		return err
	}
	// Good, run final compilation after instrumentation
	err = util.RunCmd(rp.compileArgs...)
	util.Log("RunCmd: %v (%v)", bundle.ImportPath, rp.compileArgs)
	return err
}

func Instrument() error {
	// Remove the tool itself from the command line arguments
	args := os.Args[2:]
	// Is compile command?
	if util.IsCompileCommand(strings.Join(args, " ")) {
		if config.GetConf().Verbose {
			util.Log("RunCmd: %v", args)
		}
		bundles, err := resource.LoadRuleBundles()
		if err != nil {
			err = errc.Adhere(err, "cmd", fmt.Sprintf("%v", args))
			return err
		}
		for _, bundle := range bundles {
			util.Assert(bundle.IsValid(), "sanity check")
			// Is compiling the target package?
			if matchImportPath(bundle.ImportPath, args) {
				util.Log("Apply bundle %v", bundle)
				err = compileRemix(bundle, args)
				if err != nil {
					err = errc.Adhere(err, "cmd", fmt.Sprintf("%v", args))
					err = errc.Adhere(err, "bundle", bundle.String())
					return err
				}
				return nil
			}
		}
	}
	// Not a compile command, just run it as is
	return util.RunCmd(args...)
}
