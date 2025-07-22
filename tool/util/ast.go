// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alibaba/loongsuite-go-agent/tool/errc"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const (
	IdentNil    = "nil"
	IdentTrue   = "true"
	IdentFalse  = "false"
	IdentIgnore = "_"
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
	ident.Name = IdentIgnore
	return ident
}

func IsUnusedIdent(ident *dst.Ident) bool {
	return ident.Name == IdentIgnore
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

func IsStringLit(expr dst.Expr, val string) bool {
	lit, ok := expr.(*dst.BasicLit)
	return ok &&
		lit.Kind == token.STRING &&
		lit.Value == fmt.Sprintf("%q", val)
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
	return &dst.BasicLit{Value: IdentTrue}
}

func BoolFalse() *dst.BasicLit {
	return &dst.BasicLit{Value: IdentFalse}
}

func IsInterfaceType(typ dst.Expr) bool {
	_, ok := typ.(*dst.InterfaceType)
	return ok
}

func IsEllipsis(typ dst.Expr) bool {
	_, ok := typ.(*dst.Ellipsis)
	return ok
}

func InterfaceType() *dst.InterfaceType {
	return &dst.InterfaceType{Methods: &dst.FieldList{List: nil}}
}

func ArrayType(elem dst.Expr) *dst.ArrayType {
	return &dst.ArrayType{Elt: elem}
}

func IfStmt(init dst.Stmt, cond dst.Expr,
	body, elseBody *dst.BlockStmt) *dst.IfStmt {
	return &dst.IfStmt{
		Init: dst.Clone(init).(dst.Stmt),
		Cond: dst.Clone(cond).(dst.Expr),
		Body: dst.Clone(body).(*dst.BlockStmt),
		Else: dst.Clone(elseBody).(*dst.BlockStmt),
	}
}

func IfNotNilStmt(cond dst.Expr, body, elseBody *dst.BlockStmt) *dst.IfStmt {
	var elseB dst.Stmt
	if elseBody == nil {
		elseB = nil
	} else {
		elseB = dst.Clone(elseBody).(dst.Stmt)
	}
	return &dst.IfStmt{
		Cond: &dst.BinaryExpr{
			X:  dst.Clone(cond).(dst.Expr),
			Op: token.NEQ,
			Y:  &dst.Ident{Name: IdentNil},
		},
		Body: dst.Clone(body).(*dst.BlockStmt),
		Else: elseB,
	}
}

func EmptyStmt() *dst.EmptyStmt {
	return &dst.EmptyStmt{}
}

func ExprStmt(expr dst.Expr) *dst.ExprStmt {
	return &dst.ExprStmt{X: dst.Clone(expr).(dst.Expr)}
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

func DefineStmts(lhs, rhs []dst.Expr) *dst.AssignStmt {
	return &dst.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: rhs,
	}
}

func SwitchCase(list []dst.Expr, stmts []dst.Stmt) *dst.CaseClause {
	return &dst.CaseClause{
		List: list,
		Body: stmts,
	}
}

func AddStructField(decl dst.Decl, name string, typ string) {
	gen, ok := decl.(*dst.GenDecl)
	if !ok {
		LogFatal("decl is not a GenDecl")
	}
	fd := NewField(name, Ident(typ))
	st := gen.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	st.Fields.List = append(st.Fields.List, fd)
}

func addImport(root *dst.File, paths ...string) *dst.GenDecl {
	importStmt := &dst.GenDecl{Tok: token.IMPORT}
	specs := make([]dst.Spec, 0)
	for _, path := range paths {
		spec := &dst.ImportSpec{
			Path: &dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", path),
			},
			Name: &dst.Ident{Name: IdentIgnore},
		}
		specs = append(specs, spec)
	}
	importStmt.Specs = specs
	root.Decls = append([]dst.Decl{importStmt}, root.Decls...)
	return importStmt
}

func AddImportForcely(root *dst.File, paths ...string) *dst.GenDecl {
	return addImport(root, paths...)
}

func RemoveImport(root *dst.File, path string) *dst.ImportSpec {
	for j, decl := range root.Decls {
		if genDecl, ok := decl.(*dst.GenDecl); ok &&
			genDecl.Tok == token.IMPORT {
			for i, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*dst.ImportSpec); ok {
					if importSpec.Path.Value == fmt.Sprintf("%q", path) {
						genDecl.Specs =
							append(genDecl.Specs[:i],
								genDecl.Specs[i+1:]...)
						if len(genDecl.Specs) == 0 {
							root.Decls =
								append(root.Decls[:j], root.Decls[j+1:]...)
						}
						return importSpec
					}
				}
			}
		}
	}
	return nil
}

func FindImport(root *dst.File, path string) *dst.ImportSpec {
	for _, decl := range root.Decls {
		if genDecl, ok := decl.(*dst.GenDecl); ok &&
			genDecl.Tok == token.IMPORT {
			for _, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*dst.ImportSpec); ok {
					if importSpec.Path.Value == fmt.Sprintf("%q", path) {
						return importSpec
					}
				}
			}
		}
	}
	return nil
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

func isValidRegex(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}

func MatchFuncDecl(decl dst.Decl, function string, receiverType string) bool {
	Assert(isValidRegex(function), "invalid function name pattern")

	funcDecl, ok := decl.(*dst.FuncDecl)
	if !ok {
		return false
	}
	re := regexp.MustCompile("^" + function + "$") // strict match
	if !re.MatchString(funcDecl.Name.Name) {
		return false
	}
	if receiverType != "" {
		re = regexp.MustCompile("^" + receiverType + "$") // strict match
		if !HasReceiver(funcDecl) {
			return re.MatchString("")
		}
		switch recvTypeExpr := funcDecl.Recv.List[0].Type.(type) {
		case *dst.StarExpr:
			if _, ok := recvTypeExpr.X.(*dst.Ident); !ok {
				// This is a generic type, we don't support it yet
				return false
			}
			t := "*" + recvTypeExpr.X.(*dst.Ident).Name
			return re.MatchString(t)
		case *dst.Ident:
			t := recvTypeExpr.Name
			return re.MatchString(t)
		case *dst.IndexExpr:
			// This is a generic type, we don't support it yet
			return false
		default:
			msg := fmt.Sprintf("unexpected receiver type: %T", recvTypeExpr)
			UnimplementedT(msg)
		}
	} else {
		if HasReceiver(funcDecl) {
			return false
		}
	}
	return true
}

func MatchStructDecl(decl dst.Decl, structType string) bool {
	if genDecl, ok := decl.(*dst.GenDecl); ok {
		if genDecl.Tok == token.TYPE {
			if typeSpec, ok := genDecl.Specs[0].(*dst.TypeSpec); ok {
				if typeSpec.Name.Name == structType {
					return true
				}
			}
		}
	}
	return false
}

// AST Parser
// @@ N.B. DST framework provides a series of RestoreResolvers such
// as guess.New for resolving the package name from an importPath.
// However, its strategy is simply to guess by taking last section
// of the importpath as the package name. This can lead to issues
// where package names like github.com/foo/v2 are resolved as v2,
// while in reality, they might be foo. Incorrect resolutions can
// lead to some imports that should be present being rudely removed.
// To solve this issue, we disable DST's automatic Import management
// and use plain AST manipulation to add imports.

type AstParser struct {
	fset *token.FileSet
	dec  *decorator.Decorator
}

func NewAstParser() *AstParser {
	return &AstParser{
		fset: token.NewFileSet(),
	}
}

func (ap *AstParser) FindPosition(node dst.Node) token.Position {
	astNode := ap.dec.Ast.Nodes[node]
	if astNode == nil {
		return token.Position{Filename: "", Line: -1, Column: -1} // Invalid
	}
	return ap.fset.Position(astNode.Pos())
}

// ParseSnippet parses the AST from incomplete source code snippet.
func (ap *AstParser) ParseSnippet(codeSnippet string) ([]dst.Stmt, error) {
	Assert(codeSnippet != "", "empty code snippet")
	snippet := "package main; func _() {" + codeSnippet + "}"
	file, err := decorator.ParseFile(ap.fset, "", snippet, 0)
	if err != nil {
		return nil, errc.New(errc.ErrParseCode, err.Error())
	}
	return file.Decls[0].(*dst.FuncDecl).Body.List, nil
}

// ParseSource parses the AST from complete source code.
func (ap *AstParser) ParseSource(source string) (*dst.File, error) {
	Assert(source != "", "empty source")
	ap.dec = decorator.NewDecorator(ap.fset)
	dstRoot, err := ap.dec.Parse(source)
	if err != nil {
		return nil, errc.New(errc.ErrParseCode, err.Error())
	}
	return dstRoot, nil
}

func (ap *AstParser) ParseFile(filePath string, mode parser.Mode) (*dst.File, error) {
	name := filepath.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errc.New(errc.ErrOpenFile, err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogFatal("failed to close file %s: %v", file.Name(), err)
		}
	}(file)
	astFile, err := parser.ParseFile(ap.fset, name, file, mode)
	if err != nil {
		return nil, errc.New(errc.ErrParseCode, err.Error())
	}
	ap.dec = decorator.NewDecorator(ap.fset)
	dstFile, err := ap.dec.DecorateFile(astFile)
	if err != nil {
		return nil, errc.New(errc.ErrParseCode, err.Error())
	}
	return dstFile, nil
}

func ParseAstFromFileOnlyPackage(filePath string) (*dst.File, error) {
	return NewAstParser().ParseFile(filePath, parser.PackageClauseOnly)
}

func ParseAstFromFileFast(filePath string) (*dst.File, error) {
	return NewAstParser().ParseFile(filePath, parser.SkipObjectResolution)
}

// ParseAstFromFile parses the AST from complete source file.
func ParseAstFromFile(filePath string) (*dst.File, error) {
	return NewAstParser().ParseFile(filePath, parser.ParseComments)
}

// WriteAstToFile writes the AST to source file.
func WriteAstToFile(astRoot *dst.File, filePath string) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", errc.New(errc.ErrCreateFile, err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogFatal("failed to close file %s: %v", file.Name(), err)
		}
	}(file)

	r := decorator.NewRestorer()
	err = r.Fprint(file, astRoot)
	if err != nil {
		return "", errc.New(errc.ErrParseCode, err.Error())
	}
	return file.Name(), nil
}
