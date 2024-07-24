package instrument

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	return fmt.Errorf("failed to replace compile arg %v", newArg)
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

func matchCompileImportPath(importPath string, args []string) bool {
	for _, arg := range args {
		if arg == importPath {
			return true
		}
	}
	return false
}

func Instrument() error {
	args := os.Args[2:]
	// Is compile command?
	if shared.IsCompileCommand(strings.Join(args, " ")) {
		if shared.Verbose {
			log.Printf("Compiling: %v\n", args)
		}
		bundles, err := resource.LoadRuleBundles()
		if err != nil {
			return fmt.Errorf("failed to load rule bundles: %w", err)
		}
		for _, bundle := range bundles {
			util.Assert(bundle.IsValid(), "sanity check")
			// Is compiling the target package?
			if matchCompileImportPath(bundle.ImportPath, args) {
				rp := newRuleProcessor(args, bundle.PackageName)
				err = rp.applyRules(bundle)
				if err != nil {
					return fmt.Errorf("failed to apply rules: %w", err)
				}
				log.Printf("Compiled: %v", rp.compileArgs)
				// Good, run final compilation after instrumentation
				return util.RunCmd(rp.compileArgs...)
			}
		}
	}
	// Not a compile command, just run it as is
	return util.RunCmd(args...)
}
