package echo

import (
	_ "embed"
	"github.com/dave/dst"
	"otel-auto-instrumentation/internal"
)

//go:embed setup_snippet
var setupSnippet string

//go:embed package_snippet
var packageSnippet string

func init() {

	injectors := make(map[string]internal.InjectFunc)

	injectors["echo.go"] = func(file *dst.File) error {
		for _, d := range file.Decls {
			fd, ok := d.(*dst.FuncDecl)
			if ok {
				if fd.Name.Name == "New" {
					l := len(fd.Body.List)
					ret := fd.Body.List[l-1]
					fd.Body.List[l-1] = &dst.ExprStmt{
						X: &dst.CallExpr{
							Fun: &dst.Ident{
								Name: "otelEchoNewHook",
							},
							Args: []dst.Expr{&dst.Ident{
								Name: "e",
							}},
						}}
					fd.Body.List = append(fd.Body.List, ret)
				}
			}
		}
		return nil
	}

	internal.Register(&internal.Rule{
		Name:         "echo",
		Pkg:          "github.com/labstack/echo/v4",
		SetupSnippet: setupSnippet,
		PkgSnippet:   packageSnippet,
		Injectors:    injectors,
	})
}
