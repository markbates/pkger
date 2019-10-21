package parser

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/markbates/pkger/here"
)

type visitor func(node ast.Node) (w ast.Visitor)

func (v visitor) Visit(node ast.Node) ast.Visitor {
	return v(node)
}

// inspired by https://gist.github.com/cryptix/d1b129361cea51a59af2
type file struct {
	fset     *token.FileSet
	astFile  *ast.File
	filename string
	decls    Decls
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.astFile)
}

func (f *file) find() (Decls, error) {
	if err := f.findOpenCalls(); err != nil {
		return nil, err
	}

	if err := f.findWalkCalls(); err != nil {
		return nil, err
	}

	return f.decls, nil
}

func (f *file) asValue(node ast.Node) (string, error) {
	var s string

	switch x := node.(type) {
	case *ast.BasicLit:
		s = x.Value
	case *ast.Ident:
		s = x.Name
	}

	return strconv.Unquote(s)
}

func (f *file) findOpenCalls() error {
	var err error
	f.walk(func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		exists := isPkgDot(ce.Fun, "pkger", "Open")
		if !(exists) || len(ce.Args) != 1 {
			return true
		}

		n := ce.Args[0]

		s, err := f.asValue(n)
		if err != nil {
			return false
		}

		info, err := here.Dir(filepath.Dir(f.filename))
		if err != nil {
			return false
		}

		pf := &File{
			Abs:  f.filename,
			Here: info,
		}

		decl := OpenDecl{
			file:  pf,
			pos:   n.Pos(),
			value: s,
		}

		f.decls = append(f.decls, decl)
		return true
	})
	return err
}

func (f *file) findWalkCalls() error {
	var err error
	f.walk(func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		exists := isPkgDot(ce.Fun, "pkger", "Walk")
		if !(exists) || len(ce.Args) != 2 {
			return true
		}

		n := ce.Args[0]

		s, err := f.asValue(n)
		if err != nil {
			return false
		}

		info, err := here.Dir(filepath.Dir(f.filename))
		if err != nil {
			return false
		}

		pf := &File{
			Abs:  f.filename,
			Here: info,
		}

		decl := WalkDecl{
			file:  pf,
			pos:   n.Pos(),
			value: s,
		}

		f.decls = append(f.decls, decl)
		return true
	})
	return err
}

// helpers
// =======
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// wrap a function to fulfill ast.Visitor interface
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}
