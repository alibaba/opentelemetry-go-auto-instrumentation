package instrument

import (
	_ "embed"
	"errors"
	"fmt"
	"go/token"
	"path/filepath"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

	"github.com/dave/dst"
)

// -----------------------------------------------------------------------------
// Trampoline Jump
//
// We distinguish between three types of functions: RawFunc, TrampolineFunc, and
// HookFunc. RawFunc is the original function that needs to be instrumented.
// TrampolineFunc is the function that is generated to call the onEnter and
// onExit hooks, it serves as a trampoline to the original function. HookFunc is
// the function that is called at entrypoint and exitpoint of the RawFunc. The
// so-called "Trampoline Jump" snippet is inserted at start of raw func, it is
// guaranteed to be generated within one line to avoid confusing debugging, as
// its name suggests, it jumps to the trampoline function from raw function.

const (
	TrampolineSkipName               = "skip"
	TrampolineCallContextName        = "callContext"
	TrampolineCallContextType        = "CallContext"
	TrampolineOnEnterName            = "OtelOnEnterTrampoline"
	TrampolineOnExitName             = "OtelOnExitTrampoline"
	TrampolineOnEnterNamePlaceholder = "\"OtelOnEnterNamePlaceholder\""
	TrampolineOnExitNamePlaceholder  = "\"OtelOnExitNamePlaceholder\""
)

//go:embed template.go
var trampolineTemplate string

func (rp *RuleProcessor) materializeTemplate() error {
	// Read trampoline template and materialize onEnter and onExit function
	// declarations based on that
	astRoot, err := shared.ParseAstFromSource(trampolineTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse trampoline template: %w", err)
	}
	varDecls := make([]dst.Decl, 0)
	var onEnterDecl, onExitDecl *dst.FuncDecl
	for _, node := range astRoot.Decls {
		// Materialize function declarations
		if decl, ok := node.(*dst.FuncDecl); ok {
			if decl.Name.Name == TrampolineOnEnterName {
				onEnterDecl = decl
			} else if decl.Name.Name == TrampolineOnExitName {
				onExitDecl = decl
			}
		}
		// Materialize variable declarations
		if decl, ok := node.(*dst.GenDecl); ok {
			// No further processing for variable declarations, just append them
			if decl.Tok == token.VAR {
				varDecls = append(varDecls, decl)
			}
		}
	}
	util.Assert(len(varDecls) > 0, "sanity check")
	util.Assert(onEnterDecl != nil && onExitDecl != nil, "sanity check")
	rp.onEnterHookFunc = onEnterDecl
	rp.onExitHookFunc = onExitDecl
	rp.addDecl(onEnterDecl)
	rp.addDecl(onExitDecl)
	rp.varDecls = varDecls
	return nil
}

func getNames(list *dst.FieldList) []string {
	var names []string
	for _, field := range list.List {
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
	}
	return names
}

func (rp *RuleProcessor) makeOnXName(t *api.InstFuncRule, onEnter bool) string {
	if onEnter {
		return util.GetVarNameOfFunc(t.OnEnter)
	} else {
		return util.GetVarNameOfFunc(t.OnExit)
	}
}

type ParamTrait struct {
	Index          int
	IsVaradic      bool
	IsInterfaceAny bool
}

func getHookFunc(t *api.InstFuncRule, onEnter bool) (*dst.FuncDecl, error) {
	res, err := resource.FindRuleFiles(t)
	if err != nil {
		return nil, fmt.Errorf("failed to find rule files: %w", err)
	}
	for _, file := range res {
		source, err := resource.ReadRuleFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", file, err)
		}
		astRoot, err := shared.ParseAstFromSource(source)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ast from source: %w", err)
		}
		var target *dst.FuncDecl
		if onEnter {
			target = shared.FindFuncDecl(astRoot, t.OnEnter)
		} else {
			target = shared.FindFuncDecl(astRoot, t.OnExit)
		}
		if target != nil {
			return target, nil
		}
	}
	return nil, errors.New("can not find hook func for rule")
}

func getHookParamTraits(t *api.InstFuncRule, onEnter bool) ([]ParamTrait, error) {
	target, err := getHookFunc(t, onEnter)
	if err != nil {
		return nil, fmt.Errorf("failed to get hook func: %w", err)
	}
	var attrs []ParamTrait
	// Find which parameter is type of interface{}
	for i, field := range target.Type.Params.List {
		attr := ParamTrait{Index: i}
		if _, ok := field.Type.(*dst.InterfaceType); ok {
			attr.IsInterfaceAny = true
		}
		if _, ok := field.Type.(*dst.Ellipsis); ok {
			attr.IsVaradic = true
		}
		attrs = append(attrs, attr)
	}
	return attrs, nil
}

func (rp *RuleProcessor) callOnEnterHook(t *api.InstFuncRule, traits []ParamTrait) error {
	if len(traits) != (len(rp.onEnterHookFunc.Type.Params.List) + 1 /*CallContext*/) {
		return errors.New("hook param traits mismatch on enter hook")
	}
	// Hook: 	   func onEnterFoo(callContext* CallContext, p*[]int)
	// Trampoline: func OtelOnEnterTrampoline_foo(p *[]int)
	args := []dst.Expr{dst.NewIdent(TrampolineCallContextName)}
	for idx, field := range rp.onEnterHookFunc.Type.Params.List {
		trait := traits[idx+1 /*CallContext*/]
		for _, name := range field.Names { // syntax of n1,n2 type
			if trait.IsVaradic {
				args = append(args, shared.DereferenceOf(shared.Ident(name.Name+"...")))
			} else {
				args = append(args, shared.DereferenceOf(dst.NewIdent(name.Name)))
			}
		}
	}
	// Generate onEnter call
	onEnterCall := shared.CallTo(rp.makeOnXName(t, true), args)
	insertAtBody(rp.onEnterHookFunc, shared.ExprStmt(onEnterCall),
		len(rp.onEnterHookFunc.Body.List)-1 /*before return*/)
	return nil
}

func (rp *RuleProcessor) callOnExitHook(t *api.InstFuncRule, traits []ParamTrait) error {
	if len(traits) != len(rp.onExitHookFunc.Type.Params.List) {
		return errors.New("hook param traits mismatch on exit hook")
	}
	// Hook: 	   func onExitFoo(ctx* CallContext, p*[]int)
	// Trampoline: func OtelOnExitTrampoline_foo(ctx* CallContext, p *[]int)
	var args []dst.Expr
	for idx, field := range rp.onExitHookFunc.Type.Params.List {
		if idx == 0 {
			args = append(args, dst.NewIdent(TrampolineCallContextName))
			continue
		}
		trait := traits[idx]
		for _, name := range field.Names { // syntax of n1,n2 type
			if trait.IsVaradic {
				arg := shared.DereferenceOf(shared.Ident(name.Name + "..."))
				args = append(args, arg)
			} else {
				arg := shared.DereferenceOf(dst.NewIdent(name.Name))
				args = append(args, arg)
			}
		}
	}
	// Generate onExit call
	onExitCall := shared.CallTo(rp.makeOnXName(t, false), args)
	insertAtBody(rp.onExitHookFunc, shared.ExprStmt(onExitCall),
		len(rp.onExitHookFunc.Body.List))
	return nil
}

func rectifyAnyType(paramList *dst.FieldList, traits []ParamTrait) error {
	if len(paramList.List) != len(traits) {
		return fmt.Errorf("param list length mismatch: %d vs %d",
			len(paramList.List), len(traits))
	}
	for i, field := range paramList.List {
		trait := traits[i]
		if trait.IsInterfaceAny {
			// Rectify type to "interface{}"
			field.Type = shared.InterfaceType()
		}
	}
	return nil
}

func (rp *RuleProcessor) addOnEnterHookVarDecl(t *api.InstFuncRule, traits []ParamTrait) error {
	paramTypes := rp.buildTrampolineType(true)
	addCallContext(paramTypes)
	// Hook functions may uses interface{} as parameter type, as some types of
	// raw function is not exposed, we need to use interface{} to represent them.
	err := rectifyAnyType(paramTypes, traits)
	if err != nil {
		return fmt.Errorf("failed to rectify any type: %w", err)
	}

	// Generate onEnter var decl
	varDecl := shared.NewVarDecl(rp.makeOnXName(t, true), paramTypes)
	rp.addDecl(varDecl)
	return nil
}

func (rp *RuleProcessor) addOnExitVarHookDecl(t *api.InstFuncRule, traits []ParamTrait) error {
	paramTypes := rp.buildTrampolineType(false)
	addCallContext(paramTypes)
	err := rectifyAnyType(paramTypes, traits)
	if err != nil {
		return fmt.Errorf("failed to rectify any type: %w", err)
	}

	// Generate onExit var decl
	varDecl := shared.NewVarDecl(rp.makeOnXName(t, false), paramTypes)
	rp.addDecl(varDecl)
	return nil
}

func insertAtBody(funcDecl *dst.FuncDecl, stmt dst.Stmt, index int) {
	stmts := funcDecl.Body.List
	newStmts := append(stmts[:index],
		append([]dst.Stmt{stmt}, stmts[index:]...)...)
	funcDecl.Body.List = newStmts
}

func (rp *RuleProcessor) renameFunc(t *api.InstFuncRule) {
	// Randomize trampoline function names
	rp.onEnterHookFunc.Name.Name = rp.makeFuncName(t, true)
	dst.Inspect(rp.onEnterHookFunc, func(node dst.Node) bool {
		if basicLit, ok := node.(*dst.BasicLit); ok {
			// Replace OtelOnEnterTrampolinePlaceHolder to real hook func name
			if basicLit.Value == TrampolineOnEnterNamePlaceholder {
				basicLit.Value = util.StringQuote(t.OnEnter)
			}
		}
		return true
	})
	rp.onExitHookFunc.Name.Name = rp.makeFuncName(t, false)
	dst.Inspect(rp.onExitHookFunc, func(node dst.Node) bool {
		if basicLit, ok := node.(*dst.BasicLit); ok {
			if basicLit.Value == TrampolineOnExitNamePlaceholder {
				basicLit.Value = util.StringQuote(t.OnExit)
			}
		}
		return true
	})
}

func addCallContext(list *dst.FieldList) {
	callCtx := shared.NewField(
		TrampolineCallContextName,
		shared.DereferenceOf(dst.NewIdent(TrampolineCallContextType)),
	)
	list.List = append([]*dst.Field{callCtx}, list.List...)
}

func (rp *RuleProcessor) buildTrampolineType(onEnter bool) *dst.FieldList {
	paramList := &dst.FieldList{List: []*dst.Field{}}
	if onEnter {
		if shared.HasReceiver(rp.rawFunc) {
			recvField := dst.Clone(rp.rawFunc.Recv.List[0]).(*dst.Field)
			paramList.List = append(paramList.List, recvField)
		}
		for _, field := range rp.rawFunc.Type.Params.List {
			paramField := dst.Clone(field).(*dst.Field)
			paramList.List = append(paramList.List, paramField)
		}
	} else {
		if rp.rawFunc.Type.Results != nil {
			for _, field := range rp.rawFunc.Type.Results.List {
				retField := dst.Clone(field).(*dst.Field)
				paramList.List = append(paramList.List, retField)
			}
		}
	}
	return paramList
}

func (rp *RuleProcessor) rectifyTypes() {
	onEnterHookFunc, onExitHookFunc := rp.onEnterHookFunc, rp.onExitHookFunc
	onEnterHookFunc.Type.Params = rp.buildTrampolineType(true)
	onExitHookFunc.Type.Params = rp.buildTrampolineType(false)
	candidate := []*dst.FieldList{
		onEnterHookFunc.Type.Params,
		onExitHookFunc.Type.Params,
	}
	for _, list := range candidate {
		for i := 0; i < len(list.List); i++ {
			paramField := list.List[i]
			if ft, ok := paramField.Type.(*dst.Ellipsis); ok {
				// If parameter is type of ...T, we need to convert it to *[]T
				paramField.Type = shared.DereferenceOf(shared.ArrayType(ft.Elt))
			} else {
				// Otherwise, convert it to *T as usual
				paramField.Type = shared.DereferenceOf(paramField.Type)
			}
		}
	}
	addCallContext(onExitHookFunc.Type.Params)
}

func (rp *RuleProcessor) replenishCallContext(onEnter bool) bool {
	funcDecl := rp.onEnterHookFunc
	if !onEnter {
		funcDecl = rp.onExitHookFunc
	}
	for _, stmt := range funcDecl.Body.List {
		if assignStmt, ok := stmt.(*dst.AssignStmt); ok {
			rhs := assignStmt.Rhs
			if len(rhs) == 1 {
				rhsExpr := rhs[0]
				if compositeLit, ok := rhsExpr.(*dst.CompositeLit); ok {
					elems := compositeLit.Elts
					names := getNames(funcDecl.Type.Params)
					for i, name := range names {
						if i == 0 && !onEnter {
							// SKip first callContext parameter for onExit
							continue
						}
						elems = append(elems, shared.Ident(name))
					}
					compositeLit.Elts = elems
					return true
				}
			}
		}
	}
	return false
}

func (rp *RuleProcessor) generateTrampoline(t *api.InstFuncRule, funcDecl *dst.FuncDecl) error {
	rp.rawFunc = funcDecl
	// Materialize trampoline template
	err := rp.materializeTemplate()
	if err != nil {
		return fmt.Errorf("failed to materialize template: %w", err)
	}
	// Rename onEnter and onExit trampoline function names
	rp.renameFunc(t)
	// Rectify types of onEnter and onExit trampoline funcs
	rp.rectifyTypes()
	// Generate calls to onEnter and onExit hooks
	if t.OnEnter != "" {
		traits, err := getHookParamTraits(t, true)
		if err != nil {
			return fmt.Errorf("failed to get hook param traits: %w", err)
		}
		err = rp.addOnEnterHookVarDecl(t, traits)
		if err != nil {
			return fmt.Errorf("failed to add onEnter var hook decl: %w", err)
		}
		err = rp.callOnEnterHook(t, traits)
		if err != nil {
			return fmt.Errorf("failed to call onEnter: %w", err)
		}
		if !rp.replenishCallContext(true) {
			return errors.New("failed to replenish context in onEnter hook")
		}
	}
	if t.OnExit != "" {
		traits, err := getHookParamTraits(t, false)
		if err != nil {
			return fmt.Errorf("failed to get hook param traits: %w", err)
		}
		err = rp.addOnExitVarHookDecl(t, traits)
		if err != nil {
			return fmt.Errorf("failed to add onExit var hook decl: %w", err)
		}
		err = rp.callOnExitHook(t, traits)
		if err != nil {
			return fmt.Errorf("failed to call onExit: %w", err)
		}
		if !rp.replenishCallContext(false) {
			return errors.New("failed to replenish context in onExit hook")
		}
	}
	return nil
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
	util.SaveDebugFile(pkgName+"_", path)
	return nil
}
