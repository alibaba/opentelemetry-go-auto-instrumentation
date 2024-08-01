package instrument

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

const (
	TrampolineJumpIfDesc                 = "/* TRAMPOLINE_JUMP_IF */"
	TrampolineJumpIfDescRegexp           = "/\\* TRAMPOLINE_JUMP_IF \\*/"
	TrampolineNoNewlinePlaceholder       = "/* NO_NEWWLINE_PLACEHOLDER */"
	TrampolineNoNewlinePlaceholderRegexp = "/\\* NO_NEWWLINE_PLACEHOLDER \\*/\n"
	TrampolineSemicolonPlaceholder       = "/* SEMICOLON_PLACEHOLDER */"
	TrampolineSemicolonPlaceholderRegexp = "/\\* SEMICOLON_PLACEHOLDER \\*/\n"
)

const (
	OtelAPIFile        = "otel_api.go"
	OtelTrampolineFile = "otel_trampoline.go"
)

func (rp *RuleProcessor) copyOtelApi(pkgName string) error {
	// Generate  otel_api.go at working directory
	target := filepath.Join(rp.workDir, OtelAPIFile)
	file, err := resource.CopyAPITo(target, pkgName)
	if err != nil {
		return fmt.Errorf("failed to copy otel api: %w", err)
	}
	rp.addCompileArg(file)
	return nil
}

func (rp *RuleProcessor) loadAst(filePath string) (*dst.File, error) {
	file := rp.tryRelocated(filePath)
	return shared.ParseAstFromFile(file)
}

func (rp *RuleProcessor) restoreAst(filePath string, root *dst.File) (string, error) {
	filePath = rp.tryRelocated(filePath)
	name := filepath.Base(filePath)
	newFile, err := shared.WriteAstToFile(root, filepath.Join(rp.workDir, name))
	if err != nil {
		return "", err
	}
	// Remove old filepath from compilation files and use new file
	err = rp.replaceCompileArg(newFile, func(arg string) bool {
		return arg == filePath
	})
	if err != nil {
		return "", err
	}
	return newFile, nil
}

func (rp *RuleProcessor) makeFuncName(r *api.InstFuncRule, onEnter bool) string {
	prefix := TrampolineOnExitName
	if onEnter {
		prefix = TrampolineOnEnterName
	}
	return fmt.Sprintf("%s_%s%s", prefix, r.Function, rp.rule2Suffix[r])
}

func findJumpPoint(jumpIf *dst.IfStmt) *dst.BlockStmt {
	// Multiple func rules may apply to the same function, we need to find the
	// appropriate jump point to insert trampoline jump.
	if len(jumpIf.Decs.If) == 1 && jumpIf.Decs.If[0] == TrampolineJumpIfDesc {
		// Insert trampoline jump within the else block
		elseBlock := jumpIf.Else.(*dst.BlockStmt)
		if len(elseBlock.List) > 1 {
			// One trampoline jump already exists, recursively find last one
			ifStmt, ok := elseBlock.List[len(elseBlock.List)-1].(*dst.IfStmt)
			util.Assert(ok, "unexpected statement in trampoline-jump-if")
			return findJumpPoint(ifStmt)
		} else {
			// Otherwise, this is the appropriate jump point
			return elseBlock
		}
	}
	return nil
}

func (rp *RuleProcessor) insertTJump(t *api.InstFuncRule, funcDecl *dst.FuncDecl) error {
	util.Assert(t.OnEnter != "" || t.OnExit != "", "sanity check")

	var retVals []dst.Expr // nil by default
	if retList := funcDecl.Type.Results; retList != nil {
		retVals = make([]dst.Expr, 0)
		// If return values are named, collect their names, otherwise we try to
		// name them manually for further use
		for _, field := range retList.List {
			if field.Names != nil {
				for _, name := range field.Names {
					retVals = append(retVals, dst.NewIdent(name.Name))
				}
			} else {
				retValIdent := dst.NewIdent("retVal" + util.RandomString(5))
				field.Names = []*dst.Ident{retValIdent}
				retVals = append(retVals, dst.Clone(retValIdent).(*dst.Ident))
			}
		}
	}

	// Arguments for onEnter trampoline
	args := make([]dst.Expr, 0)
	// Receiver as argument for trampoline func, if any
	if shared.HasReceiver(funcDecl) {
		if recv := funcDecl.Recv.List; recv != nil {
			receiver := recv[0].Names[0].Name
			args = append(args, shared.AddressOf(shared.Ident(receiver)))
		} else {
			util.Unimplemented()
		}
	}
	// Original function arguments as arguments for trampoline func
	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			args = append(args, shared.AddressOf(shared.Ident(name.Name)))
		}
	}

	varSuffix := util.RandomString(5)
	rp.rule2Suffix[t] = varSuffix

	// Generate the trampoline-jump-if. N.B. Note that future optimization pass
	// heavily depends on the structure of trampoline-jump-if. Any change in it
	// should be carefully examined.
	onEnterCall := shared.CallTo(rp.makeFuncName(t, true), args)
	onExitCall := shared.CallTo(rp.makeFuncName(t, false), func() []dst.Expr {
		// NB. DST framework disallows duplicated node in the
		// AST tree, we need to replicate the return values
		// as they are already used in return statement above
		clone := make([]dst.Expr, len(retVals)+1)
		clone[0] = shared.Ident(TrampolineCallContextName + varSuffix)
		for i := 1; i < len(clone); i++ {
			clone[i] = shared.AddressOf(retVals[i-1])
		}
		return clone
	}())
	tjumpInit := shared.DefineStmts(
		shared.Exprs(
			shared.Ident(TrampolineCallContextName+varSuffix),
			shared.Ident(TrampolineSkipName+varSuffix),
		),
		shared.Exprs(onEnterCall),
	)
	tjumpCond := shared.Ident(TrampolineSkipName + varSuffix)
	tjumpBody := shared.BlockStmts(
		shared.ExprStmt(onExitCall),
		shared.ReturnStmt(retVals),
	)
	tjumpElse := shared.Block(shared.DeferStmt(onExitCall))
	tjump := shared.IfStmt(tjumpInit, tjumpCond, tjumpBody, tjumpElse)
	// Add this trampoline-jump-if as optimization candidates
	rp.trampolineJumps = append(rp.trampolineJumps, &TJump{
		target: funcDecl,
		ifStmt: tjump,
		rule:   t,
	})

	// @@ Unfortunately, dst framework does not support fine-grained space control
	// i.e. there is no way to generate all above AST into one line code, we have
	// to manually format it. There we insert OtelNewlineTrampolineHolder anchor
	// so that we aware of where we should remove trailing newline.
	//	if .... { /* NO_NEWWLINE_PLACEHOLDER */
	//	... /* NO_NEWWLINE_PLACEHOLDER */
	//	} else { /* NO_NEWWLINE_PLACEHOLDER */
	//	... /* SEMICOLON_PLACEHOLDER */
	//	} /* NO_NEWWLINE_PLACEHOLDER */
	//  NEW_LINE
	{ // then block
		callExpr := tjump.Body.List[0]
		callExpr.Decorations().Start.Append(TrampolineNoNewlinePlaceholder)
		callExpr.Decorations().End.Append(TrampolineSemicolonPlaceholder)
		retStmt := tjump.Body.List[1]
		retStmt.Decorations().End.Append(TrampolineNoNewlinePlaceholder)
	}
	{ // else block
		deferStmt := tjump.Else.(*dst.BlockStmt).List[0]
		deferStmt.Decorations().Start.Append(TrampolineNoNewlinePlaceholder)
		deferStmt.Decorations().End.Append(TrampolineSemicolonPlaceholder)
		tjump.Else.Decorations().End.Append(TrampolineNoNewlinePlaceholder)
		tjump.Decs.If.Append(TrampolineJumpIfDesc) // Anchor label
	}

	// Find if there is already a trampoline-jump-if, insert new tjump if so,
	// otherwise prepend to block body
	found := false
	if len(funcDecl.Body.List) > 0 {
		firstStmt := funcDecl.Body.List[0]
		if ifStmt, ok := firstStmt.(*dst.IfStmt); ok {
			point := findJumpPoint(ifStmt)
			if point != nil {
				point.List = append(point.List, shared.EmptyStmt())
				point.List = append(point.List, tjump)
				found = true
			}
		}
	}
	if !found {
		// Outmost trampoline-jump-if may follow by user code right after else
		// block, replacing the trailing newline mandatorily breaks the code,
		// we need to insert extra new line to make replacement possible
		tjump.Decorations().After = dst.EmptyLine
		funcDecl.Body.List = append([]dst.Stmt{tjump}, funcDecl.Body.List...)
	}

	// Generate corresponding trampoline code
	err := rp.generateTrampoline(t, funcDecl)
	if err != nil {
		return err
	}
	return nil
}

func (rp *RuleProcessor) inliningTJump(filePath string) error {
	text, err := util.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Remove trailing newline
	re := regexp.MustCompile(TrampolineNoNewlinePlaceholderRegexp)
	text = re.ReplaceAllString(text, " ")
	// Replace with semicolon
	re = regexp.MustCompile(TrampolineSemicolonPlaceholderRegexp)
	text = re.ReplaceAllString(text, ";")
	// Remove trampoline jump if ideitifiers
	re = regexp.MustCompile(TrampolineJumpIfDescRegexp)
	text = re.ReplaceAllString(text, "")
	// All done, persist to file
	_, err = util.WriteStringToFile(filePath, text)
	return err
}

func (rp *RuleProcessor) insertRaw(r *api.InstFuncRule, decl *dst.FuncDecl) error {
	util.Assert(r.OnEnter != "" || r.OnExit != "", "sanity check")
	if r.OnEnter != "" {
		// Prepend raw code snippet to function body for onEnter
		onEnterSnippet, err := shared.ParseAstFromSnippet(r.OnEnter)
		if err != nil {
			return err
		}
		decl.Body.List = append(onEnterSnippet, decl.Body.List...)
	}
	if r.OnExit != "" {
		// Use defer func(){ raw_code_snippet }() for onExit
		onExitSnippet, err := shared.ParseAstFromSnippet(
			fmt.Sprintf("defer func(){ %s }()", r.OnExit),
		)
		if err != nil {
			return err
		}
		decl.Body.List = append(onExitSnippet, decl.Body.List...)
	}
	return nil
}

func nameReturnValues(funcDecl *dst.FuncDecl) {
	if funcDecl.Type.Results != nil {
		idx := 0
		for _, field := range funcDecl.Type.Results.List {
			if field.Names == nil {
				name := fmt.Sprintf("retVal%d", idx)
				field.Names = []*dst.Ident{shared.Ident(name)}
				idx++
			}
		}
	}
}

func sortFuncRules(rules []uint64) []*api.InstFuncRule {
	fnRules := make([]*api.InstFuncRule, 0)
	for _, ruleHash := range rules {
		rule := resource.FindFuncRuleByHash(ruleHash)
		fnRules = append(fnRules, rule)
	}
	sort.SliceStable(fnRules, func(i, j int) bool {
		return fnRules[i].Order < fnRules[j].Order
	})
	return fnRules
}

func (rp *RuleProcessor) writeTrampoline(pkgName string) error {
	// Prepare trampoline code header
	code := "package " + pkgName
	trampoline, err := shared.ParseAstFromSource(code)
	if err != nil {
		return fmt.Errorf("failed to parse trampoline code header: %w", err)
	}
	// One trampoline file shares common variable declarations
	trampoline.Decls = append(trampoline.Decls, rp.varDecls...)
	// Write trampoline code to file
	path := filepath.Join(rp.workDir, OtelTrampolineFile)
	trampolineFile, err := shared.WriteAstToFile(trampoline, path)
	if err != nil {
		return err
	}
	rp.addCompileArg(trampolineFile)
	// Save trampoline code for debugging
	shared.SaveDebugFile(pkgName+"_", path)
	return nil
}

func (rp *RuleProcessor) applyFuncRules(bundle *resource.RuleBundle) (err error) {
	// Copy API file to compilation working directory
	err = rp.copyOtelApi(bundle.PackageName)
	if err != nil {
		return fmt.Errorf("failed to copy otel api: %w", err)
	}

	// Applied all matched func rules, either inserting raw code or inserting
	// our trampoline calls.
	for file, fn2rules := range bundle.File2FuncRules {
		astRoot, err := rp.loadAst(file)
		if err != nil {
			return fmt.Errorf("failed to load ast from file: %w", err)
		}
		rp.target = astRoot
		rp.trampolineJumps = make([]*TJump, 0)
		for fnName, rules := range fn2rules {
			for _, decl := range astRoot.Decls {
				nameAndRecvType := strings.Split(fnName, ",")
				name := nameAndRecvType[0]
				recvType := nameAndRecvType[1]
				if resource.MatchFuncDecl(decl, name, recvType) {
					fnDecl := decl.(*dst.FuncDecl)
					// Add explicit names for return values, they can be further
					// referenced if we're willing
					nameReturnValues(fnDecl)

					// Apply all matched rules for this function
					fnRules := sortFuncRules(rules)
					for _, rule := range fnRules {
						if rule.UseRaw {
							err = rp.insertRaw(rule, fnDecl)
						} else {
							err = rp.insertTJump(rule, fnDecl)
						}
						if err != nil {
							return fmt.Errorf("failed to rewrite: %w for %v",
								err, rule)
						}
						log.Printf("Apply func rule %s\n", rule)
					}
					break
				}
			}
		}
		// Optimize generated trampoline-jump-ifs
		err = rp.optimizeTJumps()
		if err != nil {
			return fmt.Errorf("failed to optimize trampoline jumps: %w", err)
		}
		// Restore the ast to original file once all rules are applied
		filePath, err := rp.restoreAst(file, astRoot)
		if err != nil {
			return fmt.Errorf("failed to restore ast: %w", err)
		}
		// Wait, all above code snippets should be inlined to new line to
		// avoid potential misbehaviors during debugging
		err = rp.inliningTJump(filePath)
		if err != nil {
			return fmt.Errorf("failed to inline trampoline call: %w", err)
		}
		shared.SaveDebugFile("fn_", filePath)
	}

	err = rp.writeTrampoline(bundle.PackageName)
	if err != nil {
		return fmt.Errorf("failed to write trampoline: %w", err)
	}
	return nil
}
