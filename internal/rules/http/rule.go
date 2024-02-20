package http

import (
	"github.com/dave/dst"
	"go/token"
	"otel-auto-instrumentation/internal"
)

import _ "embed"

//go:embed setup_snippet
var setupSnippet string

//go:embed package_snippet
var packageSnippet string

func init() {

	injectors := make(map[string]internal.InjectFunc)

	injectors["server.go"] = func(file *dst.File) error {
		for _, d := range file.Decls {
			fd, ok := d.(*dst.FuncDecl)

			if ok {
				if fd.Name.Name == "Handle" && fd.Recv != nil {
					stmts := []dst.Stmt{
						&dst.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []dst.Expr{&dst.Ident{
								Name: "handler",
							}},
							Rhs: []dst.Expr{&dst.CallExpr{
								Fun: &dst.Ident{
									Name: "otelHandlerHook",
									Path: "",
								},
								Args: []dst.Expr{
									&dst.Ident{
										Name: "pattern",
									},
									&dst.Ident{
										Name: "handler",
									},
								},
							}},
						},
					}
					fd.Body.List = append(stmts, fd.Body.List...)
				}
			}
		}
		return nil
	}

	injectors["client.go"] = func(file *dst.File) error {
		for _, d := range file.Decls {
			fd, ok := d.(*dst.FuncDecl)
			if ok {
				if fd.Name.Name == "transport" && fd.Recv != nil {
					stmts := []dst.Stmt{
						&dst.ReturnStmt{
							Results: []dst.Expr{&dst.CallExpr{
								Fun: &dst.Ident{
									Name: "otelTransportHook",
									Path: "",
								},
								Args: []dst.Expr{
									&dst.Ident{
										Name: "c",
									},
								},
							}},
						},
					}
					fd.Body.List = stmts
				}
			}
		}
		return nil
	}

	internal.Register(&internal.Rule{
		Name:         "http",
		Pkg:          "net/http",
		SetupSnippet: setupSnippet,
		PkgSnippet:   packageSnippet,
		Injectors:    injectors,
	})
}
