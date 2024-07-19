package shared

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// AST Construction
func AddressOf(expr dst.Expr) *dst.UnaryExpr {
	return &dst.UnaryExpr{Op: token.AND, X: expr}
}

func CallTo(name string, args []dst.Expr) *dst.CallExpr {
	return &dst.CallExpr{
		Fun:  &dst.Ident{Name: name},
		Args: args,
	}
}

func MakeUnusedIdent(ident *dst.Ident) *dst.Ident {
	ident.Name = "_"
	return ident
}

func IsUnusedIdent(ident *dst.Ident) bool {
	return ident.Name == "_"
}

func Ident(name string) *dst.Ident {
	return &dst.Ident{
		Name: name,
	}
}

func Block(stmt dst.Stmt) *dst.BlockStmt {
	return &dst.BlockStmt{
		List: []dst.Stmt{
			stmt,
		},
	}
}

func BlockStmts(stmts ...dst.Stmt) *dst.BlockStmt {
	return &dst.BlockStmt{
		List: stmts,
	}
}

func Exprs(exprs ...dst.Expr) []dst.Expr {
	return exprs
}

func NewField(name string, typ dst.Expr) *dst.Field {
	newField := &dst.Field{
		Names: []*dst.Ident{dst.NewIdent(name)},
		Type:  typ,
	}
	return newField
}

func BoolTrue() *dst.BasicLit {
	return &dst.BasicLit{Value: "true"}
}

func BoolFalse() *dst.BasicLit {
	return &dst.BasicLit{Value: "false"}
}

func InterfaceType() *dst.InterfaceType {
	return &dst.InterfaceType{Methods: &dst.FieldList{List: nil}}
}

func ArrayType(elem dst.Expr) *dst.ArrayType {
	return &dst.ArrayType{Elt: elem}
}

func EmptyStmt() *dst.EmptyStmt {
	return &dst.EmptyStmt{}
}

func ExprStmt(expr dst.Expr) *dst.ExprStmt {
	return &dst.ExprStmt{X: expr}
}

func DeferStmt(call *dst.CallExpr) *dst.DeferStmt {
	return &dst.DeferStmt{Call: call}
}

func ReturnStmt(results []dst.Expr) *dst.ReturnStmt {
	return &dst.ReturnStmt{Results: results}
}

func AddStructField(decl dst.Decl, name string, typ string) {
	gen, ok := decl.(*dst.GenDecl)
	if !ok {
		log.Fatalf("decl is not a GenDecl")
	}
	fd := NewField(name, Ident(typ))
	st := gen.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	st.Fields.List = append(st.Fields.List, fd)
}

func AddImportForcely(root *dst.File, path string) {
	importStmt := &dst.GenDecl{
		Tok: token.IMPORT,
		Specs: []dst.Spec{
			&dst.ImportSpec{
				Name: &dst.Ident{Name: "_"},
				Path: &dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", path),
				},
			},
		},
	}
	root.Decls = append([]dst.Decl{importStmt}, root.Decls...)
}

func NewVarDecl(name string, paramTypes *dst.FieldList) *dst.GenDecl {
	return &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{
					{Name: name},
				},
				Type: &dst.FuncType{
					Func:   false,
					Params: paramTypes,
				},
			},
		},
	}
}

func DereferenceOf(expr dst.Expr) dst.Expr {
	return &dst.StarExpr{X: expr}
}

func HasReceiver(fn *dst.FuncDecl) bool {
	return fn.Recv != nil && len(fn.Recv.List) > 0
}

// AST utilities

func FindFuncDecl(root *dst.File, name string) *dst.FuncDecl {
	for _, decl := range root.Decls {
		if fn, ok := decl.(*dst.FuncDecl); ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}

// AST Parser

// ParseAstFromSnippet parses the AST from incomplete source code snippet.
func ParseAstFromSnippet(codeSnippnet string) ([]dst.Stmt, error) {
	fset := token.NewFileSet()
	snippet := "package main; func _() {" + codeSnippnet + "}"
	file, err := decorator.ParseFile(fset, "", snippet, 0)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*dst.FuncDecl).Body.List, nil
}

// ParseAstFromSource parses the AST from complete source code.
func ParseAstFromSource(source string) (*dst.File, error) {
	dec := decorator.NewDecorator(token.NewFileSet())
	dstRoot, err := dec.Parse(source)
	if err != nil {
		return nil, err
	}
	return dstRoot, nil
}

// ParseAstFromFile parses the AST from complete source file.
func ParseAstFromFile(filePath string) (*dst.File, error) {
	name := filepath.Base(filePath)
	fset := token.NewFileSet()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	astFile, err := parser.ParseFile(fset, name, file, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	dec := decorator.NewDecorator(fset)
	dstFile, err := dec.DecorateFile(astFile)
	if err != nil {
		return nil, err
	}
	return dstFile, nil
}

// WriteAstToFile writes the AST to source file.
func WriteAstToFile(astRoot *dst.File, filePath string) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file %s: %v", file.Name(), err)
		}
	}(file)

	r := decorator.NewRestorer()
	err = r.Fprint(file, astRoot)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
