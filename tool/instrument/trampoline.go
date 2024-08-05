package instrument

import (
	_ "embed"
	"errors"
	"fmt"
	"go/token"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
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
	TrampolineSetParamName           = "SetParam"
	TrampolineGetParamName           = "GetParam"
	TrampolineSetReturnValName       = "SetReturnVal"
	TrampolineGetReturnValName       = "GetReturnVal"
	TrampolineValIdentifier          = "val"
	TrampolineCtxIdentifier          = "c"
	TrampolineParamsIdentifier       = "Params"
	TrampolineReturnValsIdentifier   = "ReturnVals"
	TrampolineSkipName               = "skip"
	TrampolineCallContextName        = "callContext"
	TrampolineCallContextType        = "CallContext"
	TrampolineCallContextImplType    = "CallContextImpl"
	TrampolineOnEnterName            = "OtelOnEnterTrampoline"
	TrampolineOnExitName             = "OtelOnExitTrampoline"
	TrampolineOnEnterNamePlaceholder = "\"OtelOnEnterNamePlaceholder\""
	TrampolineOnExitNamePlaceholder  = "\"OtelOnExitNamePlaceholder\""
)

// @@ Modification on this trampoline template should be cautious, as it imposes
// many implicit constraints on generated code, known constraints are as follows:
// - It's performance critical, so it should be as simple as possible
// - It should not import any package because there is no guarantee that package
//   is existed in import config during the compilation, one practical approach
//   is to use function variables and setup these variables in preprocess stage
// - It should not panic as this affects user application
// - Function and variable names are coupled with the framework, any modification
//   on them should be synced with the framework

//go:embed template.go
var trampolineTemplate string

func (rp *RuleProcessor) materializeTemplate() error {
	// Read trampoline template and materialize onEnter and onExit function
	// declarations based on that
	astRoot, err := shared.ParseAstFromSource(trampolineTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse trampoline template: %w", err)
	}

	rp.varDecls = make([]dst.Decl, 0)
	rp.callCtxMethods = make([]*dst.FuncDecl, 0)
	for _, node := range astRoot.Decls {
		// Materialize function declarations
		if decl, ok := node.(*dst.FuncDecl); ok {
			if decl.Name.Name == TrampolineOnEnterName {
				rp.onEnterHookFunc = decl
				rp.addDecl(decl)
			} else if decl.Name.Name == TrampolineOnExitName {
				rp.onExitHookFunc = decl
				rp.addDecl(decl)
			} else if shared.HasReceiver(decl) {
				// We know exactly this is CallContextImpl method
				t := decl.Recv.List[0].Type.(*dst.StarExpr).X.(*dst.Ident).Name
				util.Assert(t == TrampolineCallContextImplType, "sanity check")
				rp.callCtxMethods = append(rp.callCtxMethods, decl)
				rp.addDecl(decl)
			}
		}
		// Materialize variable declarations
		if decl, ok := node.(*dst.GenDecl); ok {
			// No further processing for variable declarations, just append them
			if decl.Tok == token.VAR {
				rp.varDecls = append(rp.varDecls, decl)
			} else if decl.Tok == token.TYPE {
				rp.callCtxDecl = decl
				rp.addDecl(decl)
			}
		}
	}
	util.Assert(rp.callCtxDecl != nil, "sanity check")
	util.Assert(len(rp.varDecls) > 0, "sanity check")
	util.Assert(rp.onEnterHookFunc != nil, "sanity check")
	util.Assert(rp.onExitHookFunc != nil, "sanity check")
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

func makeOnXName(t *api.InstFuncRule, onEnter bool) string {
	if onEnter {
		return shared.GetVarNameOfFunc(t.OnEnter)
	} else {
		return shared.GetVarNameOfFunc(t.OnExit)
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
		if shared.IsInterfaceType(field.Type) {
			attr.IsInterfaceAny = true
		}
		if shared.IsEllipsis(field.Type) {
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
	// Call onEnter if it exists
	fnName := makeOnXName(t, true)
	call := shared.ExprStmt(shared.CallTo(fnName, args))
	iff := shared.IfNotNilStmt(
		dst.NewIdent(fnName),
		shared.Block(call),
		nil,
	)
	insertAt(rp.onEnterHookFunc, iff, len(rp.onEnterHookFunc.Body.List)-1)
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
	// Call onExit if it exists
	fnName := makeOnXName(t, false)
	call := shared.ExprStmt(shared.CallTo(fnName, args))
	iff := shared.IfNotNilStmt(
		dst.NewIdent(fnName),
		shared.Block(call),
		nil,
	)
	insertAtEnd(rp.onExitHookFunc, iff)
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

func (rp *RuleProcessor) addHookFuncVar(t *api.InstFuncRule, traits []ParamTrait, onEnter bool) error {
	paramTypes := rp.buildTrampolineType(onEnter)
	addCallContext(paramTypes)
	// Hook functions may uses interface{} as parameter type, as some types of
	// raw function is not exposed
	err := rectifyAnyType(paramTypes, traits)
	if err != nil {
		return fmt.Errorf("failed to rectify any type on enter: %w", err)
	}

	// Generate var decl
	varDecl := shared.NewVarDecl(makeOnXName(t, onEnter), paramTypes)
	rp.addDecl(varDecl)
	return nil
}

func insertAt(funcDecl *dst.FuncDecl, stmt dst.Stmt, index int) {
	stmts := funcDecl.Body.List
	newStmts := append(stmts[:index],
		append([]dst.Stmt{stmt}, stmts[index:]...)...)
	funcDecl.Body.List = newStmts
}

func insertAtEnd(funcDecl *dst.FuncDecl, stmt dst.Stmt) {
	insertAt(funcDecl, stmt, len(funcDecl.Body.List))
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
		dst.NewIdent(TrampolineCallContextType),
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
			paramFieldType := desugarType(paramField)
			paramField.Type = shared.DereferenceOf(paramFieldType)
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
	util.ShouldNotReachHereT("failed to replenish call context")
	return false
}

// -----------------------------------------------------------------------------
// Dynamic CallContext API Generation
//
// This is somewhat challenging, as we need to generate type-aware CallContext
// APIs, which means we need to generate a bunch of switch statements to handle
// different types of parameters. Different RawFuncs in the same package may have
// different types of parameters, all of them should have their own CallContext
// implementation, thus we need to generate a bunch of CallContextImpl{suffix}
// types and methods to handle them. The suffix is generated based on the rule
// suffix, so that we can distinguish them from each other.

// implementCallContext effectively "implements" the CallContext interface by
// renaming occurrences of CallContextImpl to CallContextImpl{suffix} in the
// trampoline template
func (rp *RuleProcessor) implementCallContext(t *api.InstFuncRule) {
	suffix := rp.rule2Suffix[t]
	structType := rp.callCtxDecl.Specs[0].(*dst.TypeSpec)
	util.Assert(structType.Name.Name == TrampolineCallContextImplType,
		"sanity check")
	structType.Name.Name += suffix             // type declaration
	for _, method := range rp.callCtxMethods { // method declaration
		method.Recv.List[0].Type.(*dst.StarExpr).X.(*dst.Ident).Name += suffix
	}
	for _, node := range []dst.Node{rp.onEnterHookFunc, rp.onExitHookFunc} {
		dst.Inspect(node, func(node dst.Node) bool {
			if ident, ok := node.(*dst.Ident); ok {
				if ident.Name == TrampolineCallContextImplType {
					ident.Name += suffix
					return false
				}
			}
			return true
		})
	}
}

func setValue(field string, idx int, typ dst.Expr) *dst.CaseClause {
	// *(c.Params[idx].(*int)) = val.(int)
	// c.Params[idx] = val iff type is interface{}
	se := shared.SelectorExpr(shared.Ident(TrampolineCtxIdentifier), field)
	ie := shared.IndexExpr(se, shared.IntLit(idx))
	te := shared.TypeAssertExpr(ie, shared.DereferenceOf(typ))
	pe := shared.ParenExpr(te)
	de := shared.DereferenceOf(pe)
	val := shared.Ident(TrampolineValIdentifier)
	assign := shared.AssignStmt(de, shared.TypeAssertExpr(val, typ))
	if shared.IsInterfaceType(typ) {
		assign = shared.AssignStmt(ie, val)
	}
	caseClause := shared.SwitchCase(
		shared.Exprs(shared.IntLit(idx)),
		shared.Stmts(assign),
	)
	return caseClause
}

func getValue(field string, idx int, typ dst.Expr) *dst.CaseClause {
	// return *(c.Params[idx].(*int))
	// return c.Params[idx] iff type is interface{}
	se := shared.SelectorExpr(shared.Ident(TrampolineCtxIdentifier), field)
	ie := shared.IndexExpr(se, shared.IntLit(idx))
	te := shared.TypeAssertExpr(ie, shared.DereferenceOf(typ))
	pe := shared.ParenExpr(te)
	de := shared.DereferenceOf(pe)
	ret := shared.ReturnStmt(shared.Exprs(de))
	if shared.IsInterfaceType(typ) {
		ret = shared.ReturnStmt(shared.Exprs(ie))
	}
	caseClause := shared.SwitchCase(
		shared.Exprs(shared.IntLit(idx)),
		shared.Stmts(ret),
	)
	return caseClause
}

func getParamClause(idx int, typ dst.Expr) *dst.CaseClause {
	return getValue(TrampolineParamsIdentifier, idx, typ)
}

func setParamClause(idx int, typ dst.Expr) *dst.CaseClause {
	return setValue(TrampolineParamsIdentifier, idx, typ)
}

func getReturnValClause(idx int, typ dst.Expr) *dst.CaseClause {
	return getValue(TrampolineReturnValsIdentifier, idx, typ)
}

func setReturnValClause(idx int, typ dst.Expr) *dst.CaseClause {
	return setValue(TrampolineReturnValsIdentifier, idx, typ)
}

// desugarType desugars parameter type to its original type, if parameter
// is type of ...T, it will be converted to []T
func desugarType(param *dst.Field) dst.Expr {
	if ft, ok := param.Type.(*dst.Ellipsis); ok {
		return shared.ArrayType(ft.Elt)
	}
	return param.Type
}

func (rp *RuleProcessor) rewriteCallContextImpl() {
	util.Assert(len(rp.callCtxMethods) > 4, "sanity check")
	var (
		methodSetParam  *dst.FuncDecl
		methodGetParam  *dst.FuncDecl
		methodGetRetVal *dst.FuncDecl
		methodSetRetVal *dst.FuncDecl
	)
	for _, decl := range rp.callCtxMethods {
		switch decl.Name.Name {
		case TrampolineSetParamName:
			methodSetParam = decl
		case TrampolineGetParamName:
			methodGetParam = decl
		case TrampolineGetReturnValName:
			methodGetRetVal = decl
		case TrampolineSetReturnValName:
			methodSetRetVal = decl
		}
	}
	// Rewrite SetParam and GetParam methods
	// Dont believe what you see in template.go, we will null out it and rewrite
	// the whole switch statement
	methodSetParamBody := methodSetParam.Body.List[0].(*dst.SwitchStmt).Body
	methodGetParamBody := methodGetParam.Body.List[0].(*dst.SwitchStmt).Body
	methodGetRetValBody := methodGetRetVal.Body.List[0].(*dst.SwitchStmt).Body
	methodSetRetValBody := methodSetRetVal.Body.List[0].(*dst.SwitchStmt).Body
	methodGetParamBody.List = nil
	methodSetParamBody.List = nil
	methodGetRetValBody.List = nil
	methodSetRetValBody.List = nil
	idx := 0
	if shared.HasReceiver(rp.rawFunc) {
		recvType := rp.rawFunc.Recv.List[0].Type
		clause := setParamClause(idx, recvType)
		methodSetParamBody.List = append(methodSetParamBody.List, clause)
		clause = getParamClause(idx, recvType)
		methodGetParamBody.List = append(methodGetParamBody.List, clause)
		idx++
	}
	for _, param := range rp.rawFunc.Type.Params.List {
		paramType := desugarType(param)
		for range param.Names {
			clause := setParamClause(idx, paramType)
			methodSetParamBody.List = append(methodSetParamBody.List, clause)
			clause = getParamClause(idx, paramType)
			methodGetParamBody.List = append(methodGetParamBody.List, clause)
			idx++
		}
	}
	// Rewrite GetReturnVal and SetReturnVal methods
	if rp.rawFunc.Type.Results != nil {
		idx = 0
		for _, retval := range rp.rawFunc.Type.Results.List {
			retType := desugarType(retval)
			for range retval.Names {
				clause := getReturnValClause(idx, retType)
				methodGetRetValBody.List = append(methodGetRetValBody.List, clause)
				clause = setReturnValClause(idx, retType)
				methodSetRetValBody.List = append(methodSetRetValBody.List, clause)
				idx++
			}
		}
	}
}

func (rp *RuleProcessor) callHookFunc(t *api.InstFuncRule, onEnter bool) error {
	traits, err := getHookParamTraits(t, onEnter)
	if err != nil {
		return fmt.Errorf("failed to get hook param traits: %w", err)
	}
	err = rp.addHookFuncVar(t, traits, onEnter)
	if err != nil {
		return fmt.Errorf("failed to add onEnter var hook decl: %w", err)
	}
	if onEnter {
		err = rp.callOnEnterHook(t, traits)
	} else {
		err = rp.callOnExitHook(t, traits)
	}
	if err != nil {
		return fmt.Errorf("failed to call onEnter: %w", err)
	}
	if !rp.replenishCallContext(onEnter) {
		return errors.New("failed to replenish context in onEnter hook")
	}
	return nil
}

func (rp *RuleProcessor) generateTrampoline(t *api.InstFuncRule, funcDecl *dst.FuncDecl) error {
	rp.rawFunc = funcDecl
	// Materialize various declarations from template file, no one wants to see
	// a bunch of manual AST code generation, isn't it?
	err := rp.materializeTemplate()
	if err != nil {
		return fmt.Errorf("failed to materialize template: %w", err)
	}
	// Implement CallContext interface
	rp.implementCallContext(t)
	// Rewrite type-aware CallContext APIs
	rp.rewriteCallContextImpl()
	// Rename trampoline functions
	rp.renameFunc(t)
	// Rectify types of trampoline functions
	rp.rectifyTypes()
	// Generate calls to hook functions
	if t.OnEnter != "" {
		err = rp.callHookFunc(t, true)
		if err != nil {
			return fmt.Errorf("failed to call onEnter: %w", err)
		}
	}
	if t.OnExit != "" {
		err = rp.callHookFunc(t, false)
		if err != nil {
			return fmt.Errorf("failed to call onExit: %w", err)
		}
	}
	return nil
}
