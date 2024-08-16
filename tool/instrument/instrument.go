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
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
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
	rule2Suffix     map[*api.InstFuncRule]string
	rawFunc         *dst.FuncDecl
	onEnterHookFunc *dst.FuncDecl
	onExitHookFunc  *dst.FuncDecl
	relocated       map[string]string
	trampolineJumps []*TJump // Optimization candidates
	callCtxDecl     *dst.GenDecl
	callCtxMethods  []*dst.FuncDecl
	assembly        string
}

const (
	OtelJumpAsm      = "otel_fancy_jump.s"
	OtelJumpSymABI   = "otel_fancy_jump_symabi"
	OtelJumpObject   = "otel_fancy_jump.o"
	GoDefaultArchive = "_pkg_.a"
)

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
		rule2Suffix: make(map[*api.InstFuncRule]string),
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

func (rp *RuleProcessor) appendCompileArg(newArg ...string) {
	rp.compileArgs = append(rp.compileArgs, newArg...)
}

func (rp *RuleProcessor) prependCompileArg(newArgs ...string) {
	for i, arg := range rp.compileArgs {
		if strings.HasSuffix(arg, "/compile") {
			// add newarg after compile
			rp.compileArgs = append(rp.compileArgs[:i+1],
				append(newArgs,
					rp.compileArgs[i+1:]...)...)
			return
		}
	}
	util.ShouldNotReachHereT("Not a compile command")
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
	return fmt.Errorf("failed to replace compile arg %v", newArg)
}

func (rp *RuleProcessor) removeCompileArg(arg string) bool {
	for i, a := range rp.compileArgs {
		if a == arg {
			rp.compileArgs = append(rp.compileArgs[:i], rp.compileArgs[i+1:]...)
			return true
		}
	}
	return false
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

// Golang has two calling convention interfaces: ABI0 and ABIInternal. To align
// with the Go runtime, we need to generate symbol ABI for the assembly code.
// See https://golang.org/design/27539-internal-abi for more details.
func (rp *RuleProcessor) genSymABI(importPath, src, dst string) (string, error) {
	util.Guarantee(importPath != "", "sanity check")
	util.Guarantee(shared.IsAsmFile(src), "sanity check")

	// Generate new symbol ABI file
	err := compileAsm(importPath, src, dst, true)
	if err != nil {
		return "", fmt.Errorf("failed to gen symabi: %w", err)
	}

	// Check if -symabis flag is already in compile arguments
	// If so, concat the old and new symabi file
	for i, arg := range rp.compileArgs {
		if arg == "-symabis" {
			newSymabi := dst
			oldSymabi := rp.compileArgs[i+1]
			err = util.ConcatFile(oldSymabi, newSymabi)
			if err != nil {
				return "", fmt.Errorf("failed to concat file: %w", err)
			}
			return oldSymabi, nil
		}
	}

	rp.prependCompileArg("-symabis", dst)
	return dst, nil
}

func compileAsm(importPath, src, dst string, genSymABI bool) error {
	util.Guarantee(importPath != "", "sanity check")
	util.Guarantee(shared.IsAsmFile(src), "sanity check")
	args := []string{shared.GoAsmTool()}
	if genSymABI {
		// Generate symbol ABI but not compile
		args = append(args, "-gensymabis")
	}
	args = append(args,
		"-o", dst,
		"-p", importPath,
		"-I", shared.GoIncludeDir(),
		src,
	)
	err := util.RunCmd(args...)
	if err != nil {
		return fmt.Errorf("failed to compile asm: %w", err)
	}
	return nil
}

func packObject(src, dst string) error {
	util.Guarantee(shared.IsObjectFile(src), "sanity check")
	util.Guarantee(shared.IsArchiveFile(dst), "sanity check")
	err := util.RunCmd(
		shared.GoPackTool(),
		"r",
		dst,
		src,
	)
	if err != nil {
		return fmt.Errorf("failed to pack object file: %w", err)
	}
	return nil
}

func compileRemix(args []string, bundle *resource.RuleBundle) error {
	// Apply rules to the target bundle
	rp := newRuleProcessor(args, bundle.PackageName)
	err := rp.applyRules(bundle)
	if err != nil {
		return fmt.Errorf("failed to apply rules: %w", err)
	}
	if rp.assembly == "" {
		// No assembly code generated, just compile as usual
		err = util.RunCmd(rp.compileArgs...)
		if err != nil {
			return fmt.Errorf("failed to compile: %w %v", err, rp.compileArgs)
		}
		return nil
	}

	// Hard case, we need to compile assembly and pack it into target archive
	// Generate symbol ABI for further compilation
	src := filepath.Join(rp.workDir, OtelJumpAsm)
	dst := filepath.Join(rp.workDir, OtelJumpSymABI)
	dst, err = rp.genSymABI(bundle.ImportPath, src, dst)
	if err != nil {
		return fmt.Errorf("failed to gen symabi: %w", err)
	}
	shared.SaveDebugFile("symabi_"+rp.packageName, dst)

	// Good, run final compilation after instrumentation
	rp.removeCompileArg("-complete")
	err = util.RunCmd(rp.compileArgs...)
	if err != nil {
		return fmt.Errorf("failed to run compile: %w %v", err, rp.compileArgs)
	}

	// Compile the generated assembly file
	src = filepath.Join(rp.workDir, OtelJumpAsm)
	dst = filepath.Join(rp.workDir, OtelJumpObject)
	err = compileAsm(bundle.ImportPath, src, dst, false)
	if err != nil {
		return fmt.Errorf("failed to compile asm: %w", err)
	}
	shared.SaveDebugFile("obj_"+rp.packageName, dst)

	// Append the assembly object file to target archive
	src = filepath.Join(rp.workDir, OtelJumpObject)
	dst = filepath.Join(rp.workDir, GoDefaultArchive)
	err = packObject(src, dst)
	if err != nil {
		return fmt.Errorf("failed to pack object file: %w", err)
	}

	// Over.
	return nil
}

func Instrument() error {
	args := os.Args[2:]
	// Is compile command?
	if shared.IsCompileCommand(strings.Join(args, " ")) {
		start := time.Now()
		if shared.Verbose {
			log.Printf("RunInit: %v\n", args)
		}
		bundles, err := resource.LoadRuleBundles()
		if err != nil {
			return fmt.Errorf("failed to load rule bundles: %w", err)
		}
		for _, bundle := range bundles {
			util.Assert(bundle.IsValid(), "sanity check")
			// Is compiling the target package?
			if matchImportPath(bundle.ImportPath, args) {
				err = compileRemix(args, bundle)
				end := time.Since(start)
				log.Printf("Instrument %v took %v\n", bundle.PackageName, end)
				if err != nil {
					return fmt.Errorf("failed to compile remix: %w", err)
				}
				return nil
			}
		}
	}
	// Not a compile command, just run it as is
	return util.RunCmd(args...)
}
