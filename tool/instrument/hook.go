// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instrument

import (
	"embed"
	"errors"
	"fmt"
	"go/token"
	"runtime"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/dave/dst"
)

// == Code conjured by yyang, Aug. 2024 ==

//------------------------------------------------------------------------------
// Fancy Jump To Hook
//
// Well, this is the fanciest part of the trampoline. We want the hook function
// to continue linking and compiling even when a suitable relocation target is
// not found, and to correctly link when they do exist. This is not really easily
// achievable through regular golang code. To accomplish this, we declare some
// uintptr variables to store the address of the hook function, and then use a
// fancy jump function implemented in assembly code to treat this variable as a
// function pointer and jump to it to execute.
// Since the hook function has parameters, in order to not break golang's ABI,
// we need the golang compiler to correctly pass these parameters to the fancy
// jump, which transparently passes them again to the final hook function and
// executes. During linking, if the hook function does not exist, the initial
// value of these variables is 0, and the fancy jump will not execute. If the
// hook function does exist, they are initialized to the function address of the
// hook function, and the fancy jump will correctly jump to them and execute.
// The overall process is as follows:
//
// var HookFunc uintptr
//
// func FancyJumpHookFunc(callContext* CallContext, arg1 int)
//	lea HookFunc, rax ;; Load address of hook function
//	mov (rax), rbx    ;; Load value of hook function
//	cmp rbx, 0        ;; Check if it correctly linked
//	je invalid        ;; Jump to invalid if not linked
//	jmp rax           ;; Jump to hook function
// invalid:
//	ret               ;; Return if not linked
//
// func OtelOnEnterTrampoline_foo() {
//   ...
//	 FancyJumpHookFunc(&callContext, 1)
//}

const FancyJumpHookFuncVar = "HookFuncVar"
const FancyJump = "FancyJump"

func (rp *RuleProcessor) newFancyJump(t *api.InstFuncRule, onEnter bool) error {
	asm, err := selectAsmByArch()
	if err != nil {
		return fmt.Errorf("failed to select asm by arch: %w", err)
	}
	asm = strings.ReplaceAll(asm, FancyJumpHookFuncVar, makeOnXName(t, onEnter))
	asm = strings.ReplaceAll(asm, FancyJump, jumpTarget(t, onEnter))
	asm += "\n"
	rp.assembly += asm
	return nil
}

//go:embed fancy_jump_*.s
var fancyJumpAsmFS embed.FS

func selectAsmByArch() (string, error) {
	arch := runtime.GOARCH
	bs, err := fancyJumpAsmFS.ReadFile("fancy_jump_" + arch + ".s")
	if err != nil {
		return "", fmt.Errorf("failed to read fancy jump asm file: %w arch: %s",
			err, arch)
	}
	content := string(bs)
	content = shared.RemoveGoComment(content)
	return content, nil
}

func jumpTarget(t *api.InstFuncRule, onEnter bool) string {
	return FancyJump + makeOnXName(t, onEnter)
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
	fnName := jumpTarget(t, true)
	call := shared.ExprStmt(shared.CallTo(fnName, args))
	insertAt(rp.onEnterHookFunc, call, len(rp.onEnterHookFunc.Body.List)-1)
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
	fnName := jumpTarget(t, false)
	call := shared.ExprStmt(shared.CallTo(fnName, args))
	insertAtEnd(rp.onExitHookFunc, call)
	return nil
}

func (rp *RuleProcessor) addHookFunc(t *api.InstFuncRule,
	traits []ParamTrait, onEnter bool) error {
	paramTypes := rp.buildTrampolineType(onEnter)
	addCallContext(paramTypes)
	// Hook functions may uses interface{} as parameter type, as some types of
	// raw function is not exposed
	err := rectifyAnyType(paramTypes, traits)
	if err != nil {
		return fmt.Errorf("failed to rectify any type on enter: %w", err)
	}

	// Generate uintptr hook point
	varDecl := &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{
					{Name: makeOnXName(t, onEnter)},
				},
				Type: &dst.Ident{Name: "uintptr"},
			},
		},
	}
	rp.addDecl(varDecl)
	// Generate fancy jump hook function declaration
	funcDecl := &dst.FuncDecl{
		Name: dst.NewIdent(jumpTarget(t, onEnter)),
		Type: &dst.FuncType{
			Params: paramTypes,
		},
	}
	rp.addDecl(funcDecl)
	// Generate assembly of fancy jump, which finally jumps to hook function
	rp.newFancyJump(t, onEnter)
	return nil
}

func (rp *RuleProcessor) callHookFunc(t *api.InstFuncRule, onEnter bool) error {
	traits, err := getHookParamTraits(t, onEnter)
	if err != nil {
		return fmt.Errorf("failed to get hook param traits: %w", err)
	}
	err = rp.addHookFunc(t, traits, onEnter)
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
