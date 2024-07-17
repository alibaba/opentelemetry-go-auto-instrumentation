package instrument

import (
	"fmt"
	"log"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

// == Code conjured by yyang, Jun. 2024 ==

// -----------------------------------------------------------------------------
// Trampoline Optimization
//
// Since trampoline-jump-if and trampoline functions are performance-critical,
// we are trying to optimize them as much as possible. The standard form of
// trampoline-jump-if looks like
//
//	if ctx, skip := otel_trampoline_onenter(&arg); skip {
//	    otel_trampoline_onexit(ctx, &retval)
//	    return ...
//	} else {
//	    defer otel_trampoline_onexit(ctx, &retval)
//	    ...
//	}
//
// The obvious optimization opportunities are cases when onEnter or onExit hooks
// are not present. For the latter case, we can replace the defer statement to
// empty statement, you might argue that we can remove the whole else block, but
// there might be more than one trampoline-jump-if in the same function, they are
// nested in the else block, i.e.
//
//	if ctx, skip := otel_trampoline_onenter(&arg); skip {
//	    otel_trampoline_onexit(ctx, &retval)
//	    return ...
//	} else {
//	    ;
//	    ...
//	}
//
// For the former case, it's a bit more complicated. We need to manually construct
// CallContext on the fly and pass it to onExit trampoline defer call and rewrite
// the whole condition to always false. The corresponding code snippet is
//
//	if false {
//	    ;
//	} else {
//	    defer otel_trampoline_onexit(&CallContext{...}, &retval)
//	    ...
//	}
//
// The if skeleton should be kept as is, otherwise inlining of trampoline-jump-if
// will not work. During compiling, the dce and sccp passes will remove the whole
// then block. That's not the whole story. We can further optimize the tjump iff
// the onEnter hook does not use SkipCall. In this case, we can rewrite condition
// of trampoline-jump-if to always false, remove return statement in then block,
// they are memory-aware and may generate memory SSA values during compilation.
//
//	if ctx,_ := otel_trampoline_onenter(&arg); false {
//	    ;
//	} else {
//	    defer otel_trampoline_onexit(ctx, &retval)
//	    ...
//	}
//
// The compiler responsible for hoisting the initialization statement out of the
// if skeleton, and the dce and sccp passes will remove the whole then block. All
// these trampoline functions looks as if they are executed sequentially.
// Note that this optimization pass is fraigle as it really heavily depends on
// the structure of trampoline-jump-if and trampoline functions. Any change in
// tjump should be carefully examined.

// TJump describes a trampoline-jump-if optimization candidate
type TJump struct {
	target *dst.FuncDecl     // Target function we are hooking on
	ifStmt *dst.IfStmt       // Trampoline-jump-if statement
	rule   *api.InstFuncRule // Rule associated with the trampoline-jump-if
}

func newDecoratedEmptyStmt() *dst.EmptyStmt {
	emptyStmt := shared.EmptyStmt()
	emptyStmt.Decorations().Start.Append(TrampolineNoNewlinePlaceholder)
	emptyStmt.Decorations().End.Append(TrampolineSemicolonPlaceholder)
	return emptyStmt
}

func mustTJump(ifStmt *dst.IfStmt) {
	util.Assert(len(ifStmt.Decs.If) == 1, "must be a trampoline-jump-if")
	desc := ifStmt.Decs.If[0]
	util.Assert(desc == TrampolineJumpIfDesc, "must be a trampoline-jump-if")
}

func (rp *RuleProcessor) removeOnExitTrampolineCall(tjump *TJump) error {
	ifStmt := tjump.ifStmt
	elseBlock := ifStmt.Else.(*dst.BlockStmt)
	for i, stmt := range elseBlock.List {
		if _, ok := stmt.(*dst.DeferStmt); ok {
			// Replace defer statement with an decorated empty statement to make
			// trampoline-jump-if inlining work
			elseBlock.List[i] = newDecoratedEmptyStmt()
			if shared.Verbose {
				log.Printf("Optimize tjump branch in %s",
					tjump.target.Name.Name)
			}
			break
		} else if _, ok := stmt.(*dst.IfStmt); ok {
			// Expected statement type and do nothing
		} else {
			// Unexpected statement type
			util.ShouldNotReachHereT("unexpected statement type")
		}
	}
	return nil
}

func replenishCallContextLiteral(tjump *TJump, expr dst.Expr) {
	rawFunc := tjump.target
	names := make([]dst.Expr, 0)
	for _, name := range getNames(rawFunc.Type.Params) {
		names = append(names, shared.AddressOf(shared.Ident(name)))
	}
	elems := expr.(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts
	paramLiteral := elems[0].(*dst.KeyValueExpr).Value.(*dst.CompositeLit)
	paramLiteral.Elts = names
}

func newCallContext(tjump *TJump) (dst.Expr, error) {
	// One line please, otherwise debugging line number will be a nightmare
	const newCallContext = `&CallContext{Params:[]interface{}{},ReturnVals:[]interface{}{},SkipCall:false,}`
	astRoot, err := shared.ParseAstFromSnippet(newCallContext)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new CallContext: %w", err)
	}
	ctxExpr := astRoot[0].(*dst.ExprStmt).X
	// Replenish call context by passing addresses of all arguments
	replenishCallContextLiteral(tjump, ctxExpr)
	return ctxExpr, nil
}

func (rp *RuleProcessor) removeOnEnterTrampolineCall(tjump *TJump) error {
	// Construct CallContext on the fly and pass to onExit trampoline defer call
	callContextExpr, err := newCallContext(tjump)
	if err != nil {
		return fmt.Errorf("failed to construct CallContext: %w", err)
	}
	// Find defer call to onExit and replace its call context with new one
	found := false
	for _, stmt := range tjump.ifStmt.Else.(*dst.BlockStmt).List {
		// Replace call context argument of defer statement to structure literal
		if deferStmt, ok := stmt.(*dst.DeferStmt); ok {
			args := deferStmt.Call.Args
			util.Assert(len(args) >= 1, "must have at least one argument")
			args[0] = callContextExpr
			found = true
			break
		}
	}
	util.Assert(found, "defer statement not found")
	// Rewrite condition of trampoline-jump-if to always false and null out its
	// initialization statement and then block
	tjump.ifStmt.Init = nil
	tjump.ifStmt.Cond = shared.BoolFalse()
	tjump.ifStmt.Body = shared.Block(newDecoratedEmptyStmt())
	if shared.Verbose {
		log.Printf("Optimize tjump branch in %s", tjump.target.Name.Name)
	}
	// Remove generated onEnter trampoline function
	removed := rp.removeDeclWhen(func(d dst.Decl) bool {
		if funcDecl, ok := d.(*dst.FuncDecl); ok {
			return funcDecl.Name.Name == rp.makeFuncName(tjump.rule, true)

		}
		return false
	})
	if removed == nil {
		return fmt.Errorf("failed to remove onEnter trampoline function")
	}
	return nil
}

func flattenTJump(tjump *TJump, removedOnExit bool) {
	ifStmt := tjump.ifStmt
	initStmt := ifStmt.Init.(*dst.AssignStmt)
	util.Assert(len(initStmt.Lhs) == 2, "must be")

	ifStmt.Cond = shared.BoolFalse()
	ifStmt.Body = shared.Block(newDecoratedEmptyStmt())

	if removedOnExit {
		// We removed the last reference to call context after nulling out body
		// block, at this point, all lhs are unused, replace assignment to simple
		// function call
		ifStmt.Init = shared.ExprStmt(initStmt.Rhs[0])
		// TODO: Remove onExit declaration
	} else {
		// Otherwise, mark skipCall identifier as unused
		skipCallIdent := initStmt.Lhs[1].(*dst.Ident)
		shared.MakeUnusedIdent(skipCallIdent)
	}
	if shared.Verbose {
		log.Printf("Optimize skipCall in %s", tjump.target.Name.Name)
	}
}

func (rp *RuleProcessor) optimizeTJumps() (err error) {
	for _, tjump := range rp.trampolineJumps {
		mustTJump(tjump.ifStmt)
		// No onExit hook present? Simply remove defer call to onExit trampoline.
		// Why we dont remove the whole else block of trampoline-jump-if? Well,
		// because there might be more than one trampoline-jump-if in the same
		// function, they are nested in the else block. See findJumpPoint for
		// more details.
		rule := tjump.rule
		removedOnExit := false
		if rule.OnExit == "" {
			err = rp.removeOnExitTrampolineCall(tjump)
			if err != nil {
				return fmt.Errorf("failed to optimize tjump: %w", err)
			}
			removedOnExit = true
		}
		// No onEnter hook present? Construct CallContext on the fly and pass it
		// to onExit trampoline defer call and rewrite the whole condition to
		// always false, then null out its initialization statement.
		if rule.OnEnter == "" {
			err = rp.removeOnEnterTrampolineCall(tjump)
			if err != nil {
				return fmt.Errorf("failed to optimize tjump: %w", err)
			}
		}

		// No SkipCall used in onEnter hook? Rewrite cond of trampoline-jump-if
		// to always false, and remove return statement in then block, they are
		// memory aware and may generate memory SSA values during compilation.
		// This further simplifies the trampoline-jump-if and gives more chances
		// for optimization passes to kick in.
		if rule.OnEnter != "" {
			onEnterHook, err := getHookFunc(rule, true)
			if err != nil {
				return fmt.Errorf("failed to get hook: %w", err)
			}
			foundPoison := false
			const poison = "SkipCall"
			// FIXME: We should traverse the call graph to find all possible
			// usage of SkipCall, but for now, we just check the onEnter hook
			// function body.
			dst.Inspect(onEnterHook, func(node dst.Node) bool {
				if ident, ok := node.(*dst.Ident); ok {
					if strings.Contains(ident.Name, poison) {
						foundPoison = true
						return false
					}
				}
				return true
			})
			if !foundPoison {
				flattenTJump(tjump, removedOnExit)
			}
		}
	}
	return nil
}
