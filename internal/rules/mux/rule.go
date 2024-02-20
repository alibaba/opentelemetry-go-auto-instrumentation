package mux

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

	injectors["mux.go"] = func(file *dst.File) error {
		for _, d := range file.Decls {
			fd, ok := d.(*dst.FuncDecl)

			if ok {
				if fd.Name.Name == "NewRouter" {
					len := len(fd.Body.List)
					ret := fd.Body.List[len-1].(*dst.ReturnStmt)
					ret.Results[0] = &dst.CallExpr{
						Fun: &dst.Ident{
							Name: "otelMuxNewRouterHook",
						},
						Args: []dst.Expr{
							ret.Results[0],
						},
					}
				}
			}
		}
		return nil
	}

	internal.Register(&internal.Rule{
		Name:         "mux",
		Pkg:          "github.com/gorilla/mux",
		SetupSnippet: setupSnippet,
		PkgSnippet:   packageSnippet,
		Injectors:    injectors,
	})
}
