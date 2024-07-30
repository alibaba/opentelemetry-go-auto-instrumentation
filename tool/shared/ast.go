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
	return &dst.UnaryExpr{Op: token.AND, X: dst.Clone(expr).(dst.Expr)}
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

func StringLit(value string) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.STRING,
		Value: fmt.Sprintf("%q", value),
	}
}

func IntLit(value int) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.INT,
		Value: fmt.Sprintf("%d", value),
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

func Stmts(stmts ...dst.Stmt) []dst.Stmt {
	return stmts
}

func SelectorExpr(x dst.Expr, sel string) *dst.SelectorExpr {
	return &dst.SelectorExpr{
		X:   dst.Clone(x).(dst.Expr),
		Sel: Ident(sel),
	}
}

func IndexExpr(x dst.Expr, index dst.Expr) *dst.IndexExpr {
	return &dst.IndexExpr{
		X:     dst.Clone(x).(dst.Expr),
		Index: dst.Clone(index).(dst.Expr),
	}
}

func TypeAssertExpr(x dst.Expr, typ dst.Expr) *dst.TypeAssertExpr {
	return &dst.TypeAssertExpr{
		X:    x,
		Type: dst.Clone(typ).(dst.Expr),
	}
}

func ParenExpr(x dst.Expr) *dst.ParenExpr {
	return &dst.ParenExpr{
		X: dst.Clone(x).(dst.Expr),
	}
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
	return &dst.DeferStmt{Call: dst.Clone(call).(*dst.CallExpr)}
}

func ReturnStmt(results []dst.Expr) *dst.ReturnStmt {
	return &dst.ReturnStmt{Results: results}
}

func AssignStmt(lhs, rhs dst.Expr) *dst.AssignStmt {
	return &dst.AssignStmt{
		Lhs: []dst.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []dst.Expr{rhs},
	}
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
