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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
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
	target          *dst.File // The target file to be instrumented
	compileArgs     []string
	rule2Suffix     map[*resource.InstFuncRule]string
	rawFunc         *dst.FuncDecl
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

func (rp *RuleProcessor) replaceCompileArg(newArg string, pred func(string) bool) error {
	for i, arg := range rp.compileArgs {
		if pred(arg) {
			rp.compileArgs[i] = newArg
			// Relocate the replaced file to new target, any rules targeting the
			// replaced file should be updated to target the new file as well
			rp.setRelocated(arg, newArg)
			return nil
		}
	}
	return errors.New("no matching compile arg found")
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
	dest = shared.GetInstrumentLogPath(dest)
	err := os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	if err != nil { // error is tolerable here
		log.Printf("failed to create debug file directory %s: %v", dest, err)
		return
	}
	err = util.CopyFile(path, dest)
	if err != nil { // error is tolerable here
		log.Printf("failed to save debug file %s: %v", dest, err)
	}
}

func (rp *RuleProcessor) applyRules(bundle *resource.RuleBundle) (err error) {
	// Apply file instrument rules first
	err = rp.applyFileRules(bundle)
	if err != nil {
		return fmt.Errorf("failed to apply file rules: %w", err)
	}

	err = rp.applyStructRules(bundle)
	if err != nil {
		return fmt.Errorf("failed to apply struct rules: %w", err)
	}

	err = rp.applyFuncRules(bundle)
	if err != nil {
		return fmt.Errorf("failed to apply function rules: %w %v", err, bundle)
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

// guaranteeVersion makes sure the rule bundle is still valid
func guaranteeVersion(bundle *resource.RuleBundle, candidates []string) error {
	for _, candidate := range candidates {
		// It's not a go file, ignore silently
		if !shared.IsGoFile(candidate) {
			continue
		}
		version := shared.ExtractVersion(candidate)
		for _, funcRules := range bundle.File2FuncRules {
			for _, rules := range funcRules {
				for _, rule := range rules {
					matched, err := shared.MatchVersion(version, rule.GetVersion())
					if err != nil || !matched {
						return fmt.Errorf("failed to match version %v", err)
					}
				}
			}
		}
		for _, fileRule := range bundle.FileRules {
			matched, err := shared.MatchVersion(version, fileRule.GetVersion())
			if err != nil || !matched {
				return fmt.Errorf("failed to match version %v", err)
			}
		}
		for _, structRules := range bundle.File2StructRules {
			for _, rules := range structRules {
				for _, rule := range rules {
					matched, err := shared.MatchVersion(version, rule.GetVersion())
					if err != nil || !matched {
						return fmt.Errorf("failed to match version %v", err)
					}
				}
			}
		}
		// Good, the bundle is still valid, we not need to check all files
		// in the package as they are mostly the same version
		break
	}
	return nil
}

func compileRemix(bundle *resource.RuleBundle, args []string) error {
	start := time.Now()
	guaranteeVersion(bundle, args)
	rp := newRuleProcessor(args, bundle.PackageName)
	err := rp.applyRules(bundle)
	if err != nil {
		return fmt.Errorf("failed to apply rules: %w", err)
	}
	// Good, run final compilation after instrumentation
	err = util.RunCmd(rp.compileArgs...)
	if shared.Verbose {
		log.Printf("RunCmd: %v (%v)\n",
			rp.compileArgs, time.Since(start))
	} else {
		log.Printf("RunCmd: %v (%v)\n",
			bundle.ImportPath, time.Since(start))
	}
	return err
}

func Instrument() error {
	args := os.Args[2:]
	// Is compile command?
	if shared.IsCompileCommand(strings.Join(args, " ")) {
		if shared.Verbose {
			log.Printf("RunCmd: %v\n", args)
		}
		bundles, err := resource.LoadRuleBundles()
		if err != nil {
			return fmt.Errorf("failed to load rule bundles: %w", err)
		}
		for _, bundle := range bundles {
			util.Assert(bundle.IsValid(), "sanity check")
			// Is compiling the target package?
			if matchImportPath(bundle.ImportPath, args) {
				err = compileRemix(bundle, args)

				return err
			}
		}
	}
	// Not a compile command, just run it as is
	return util.RunCmd(args...)
}
