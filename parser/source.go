package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/markbates/pkger/here"
)

type Source struct {
	Abs  string // full path on disk to file
	Path here.Path
	Here here.Info
}

type ParsedSource struct {
	Source
	FileSet *token.FileSet
	Ast     *ast.File
	decls   map[string]Decls
	once    sync.Once
	err     error
}

func (p *ParsedSource) Parse() error {
	(&p.once).Do(func() {
		p.err = p.parse()
	})
	return p.err
}

func (p *ParsedSource) valueIdent(node *ast.Ident) (s string) {
	s = node.Name
	if node.Obj.Kind != ast.Con {
		return
	}
	// As per ast package a Con object is always a *ValueSpec,
	// but double-checking to avoid panics
	if x, ok := node.Obj.Decl.(*ast.ValueSpec); ok {
		// The const var can be defined inline with other vars,
		// as in `const a, b = "a", "b"`.
		for i, v := range x.Names {
			if v.Name == node.Name {
				s = p.valueNode(x.Values[i])
				break
			}
		}
	}
	return
}

func (p *ParsedSource) valueNode(node ast.Node) string {
	var s string
	switch x := node.(type) {
	case *ast.BasicLit:
		s = x.Value
	case *ast.Ident:
		s = p.valueIdent(x)
	}
	return s
}

func (p *ParsedSource) value(node ast.Node) (string, error) {
	s := p.valueNode(node)
	return strconv.Unquote(s)
}

func (p *ParsedSource) parse() error {
	p.decls = map[string]Decls{}
	var fn walker = func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkg, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if pkg.Name != "pkger" {
			return true
		}

		var fn func(f File, pos token.Position, value string) Decl

		name := sel.Sel.Name
		switch name {
		case "MkdirAll":
			fn = func(f File, pos token.Position, value string) Decl {
				return MkdirAllDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		case "Create":
			fn = func(f File, pos token.Position, value string) Decl {
				return CreateDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		case "Include":
			fn = func(f File, pos token.Position, value string) Decl {
				return IncludeDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		case "Stat":
			fn = func(f File, pos token.Position, value string) Decl {
				return StatDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		case "Open":
			fn = func(f File, pos token.Position, value string) Decl {
				return OpenDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		case "Dir":
			fn = func(f File, pos token.Position, value string) Decl {
				return HTTPDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		case "Walk":
			fn = func(f File, pos token.Position, value string) Decl {
				return WalkDecl{
					file:  &f,
					pos:   pos,
					value: value,
				}
			}
		default:
			return true
		}

		if len(ce.Args) < 1 {
			p.err = fmt.Errorf("declarations require at least one argument")
			return false
		}

		n := ce.Args[0]
		val, err := p.value(n)
		if err != nil {
			p.err = fmt.Errorf("%s: %s", err, n)
			return false
		}

		info, err := here.Dir(filepath.Dir(p.Abs))
		if err != nil {
			p.err = fmt.Errorf("%s: %s", err, p.Abs)
			return false
		}

		pt, err := info.Parse(val)
		if err != nil {
			p.err = fmt.Errorf("%s: %s", err, p.Abs)
			return false
		}

		if pt.Pkg != info.Module.Path {
			info, err = here.Package(pt.Pkg)
			if err != nil {
				p.err = fmt.Errorf("%s: %s", err, p.Abs)
				return false
			}
		}

		f := File{
			Abs:  filepath.Join(info.Module.Dir, pt.Name),
			Here: info,
			Path: pt,
		}

		p.decls[name] = append(p.decls[name], fn(f, p.FileSet.Position(n.Pos()), val))
		return true
	}
	ast.Walk(fn, p.Ast)
	return nil
}

func (p *ParsedSource) DeclsMap() (map[string]Decls, error) {
	err := p.Parse()
	return p.decls, err
}

// wrap a function to fulfill ast.Visitor interface
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}
