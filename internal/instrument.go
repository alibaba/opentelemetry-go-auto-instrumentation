package internal

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type InjectFunc func(*dst.File) error

type Rule struct {
	Name         string
	Pkg          string
	SetupSnippet string
	PkgSnippet   string
	Injectors    map[string]InjectFunc
}

var rules = make(map[string]*Rule)

func Register(r *Rule) {
	rules[r.Pkg] = r
}

func apply(r *Rule, args []string, wd string) ([]string, error) {
	var err error
	if len(r.Injectors) > 0 {
		for i, arg := range args {
			if strings.HasSuffix(arg, ".go") {
				name := filepath.Base(arg)
				injector, exist := r.Injectors[name]

				if exist {
					fset := token.NewFileSet()
					file, _ := os.Open(arg)
					astFile, _ := parser.ParseFile(fset, name, file, parser.ParseComments)
					dec := decorator.NewDecorator(fset)
					dstFile, _ := dec.DecorateFile(astFile)
					_ = injector(dstFile)
					r := decorator.NewRestorer()
					dest := filepath.Join(wd, name)
					args[i] = dest
					file, _ = os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
					err = r.Fprint(file, dstFile)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if len(r.PkgSnippet) > 0 {
		s := filepath.Join(wd, "otel_snippet.go")
		f, _ := os.OpenFile(s, os.O_CREATE|os.O_WRONLY, 0644)
		_, err = f.WriteString(r.PkgSnippet)
		if err != nil {
			return nil, err
		}
		args = append(args, s)
	}
	return args, nil
}
