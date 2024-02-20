package gin

import (
	"github.com/dave/dst"
	"otel-auto-instrumentation/internal"
)

import _ "embed"

//go:embed setup_snippet
var setupSnippet string

//go:embed package_snippet
var packageSnippet string

func init() {
	injectors := make(map[string]internal.InjectFunc)

	injectors["gin.go"] = func(file *dst.File) error {
		for _, d := range file.Decls {
			fd, ok := d.(*dst.FuncDecl)

			if ok {
				if fd.Name.Name == "New" {
					len := len(fd.Body.List)
					ret := fd.Body.List[len-1].(*dst.ReturnStmt)
					ret.Results[0] = &dst.CallExpr{
						Fun: &dst.Ident{
							Name: "otelGinNewHook",
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
		Name:         "gin",
		Pkg:          "github.com/gin-gonic/gin",
		SetupSnippet: setupSnippet,
		PkgSnippet:   packageSnippet,
		Injectors:    injectors,
	})
}
