package grpc

import (
	_ "embed"
	"github.com/dave/dst"
	"go/token"
	"otel-auto-instrumentation/internal"
)

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
				if fd.Name.Name == "NewServer" {
					stmts := []dst.Stmt{
						&dst.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []dst.Expr{&dst.Ident{
								Name: "opt",
							}},
							Rhs: []dst.Expr{&dst.CallExpr{
								Fun: &dst.Ident{
									Name: "otelNewServerHook",
									Path: "",
								},
								Args: []dst.Expr{
									&dst.Ident{
										Name: "opt",
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

	injectors["clientconn.go"] = func(file *dst.File) error {
		for _, d := range file.Decls {
			fd, ok := d.(*dst.FuncDecl)

			if ok {
				if fd.Name.Name == "DialContext" {
					stmts := []dst.Stmt{
						&dst.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []dst.Expr{&dst.Ident{
								Name: "opts",
							}},
							Rhs: []dst.Expr{&dst.CallExpr{
								Fun: &dst.Ident{
									Name: "otelDialContextHook",
									Path: "",
								},
								Args: []dst.Expr{
									&dst.Ident{
										Name: "opts",
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

	internal.Register(&internal.Rule{
		Name:         "grpc",
		Pkg:          "google.golang.org/grpc",
		SetupSnippet: setupSnippet,
		PkgSnippet:   packageSnippet,
		Injectors:    injectors,
	})
}
