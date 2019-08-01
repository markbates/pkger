package parser

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/markbates/pkger/paths"
)

type Visitor struct {
	File   string
	Found  map[paths.Path]bool
	errors []error
}

func NewVisitor(p string) (*Visitor, error) {
	return &Visitor{
		File:  p,
		Found: map[paths.Path]bool{},
	}, nil
}

func (v *Visitor) Run() ([]paths.Path, error) {
	pf, err := ParseFile(v.File)
	if err != nil {
		return nil, err
	}

	ast.Walk(v, pf.Ast)

	var found []paths.Path

	for k := range v.Found {
		found = append(found, k)
	}

	return found, nil
}

func (v *Visitor) addPath(p string) error {
	p, _ = strconv.Unquote(p)
	pt, err := paths.Parse(p)
	if err != nil {
		return err
	}

	v.Found[pt] = true

	return nil
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	if err := v.eval(node); err != nil {
		v.errors = append(v.errors, err)
	}

	return v
}

func (v *Visitor) eval(node ast.Node) error {
	switch t := node.(type) {
	case *ast.CallExpr:
		return v.evalExpr(t)
	case *ast.Ident:
		return v.evalIdent(t)
	case *ast.GenDecl:
		for _, n := range t.Specs {
			if err := v.eval(n); err != nil {
				return err
			}
		}
	case *ast.FuncDecl:
		if t.Body == nil {
			return nil
		}
		for _, b := range t.Body.List {
			if err := v.evalStmt(b); err != nil {
				return err
			}
		}
		return nil
	case *ast.ValueSpec:
		for _, e := range t.Values {
			if err := v.evalExpr(e); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visitor) evalStmt(stmt ast.Stmt) error {
	switch t := stmt.(type) {
	case *ast.ExprStmt:
		return v.evalExpr(t.X)
	case *ast.AssignStmt:
		for _, e := range t.Rhs {
			if err := v.evalArgs(e); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visitor) evalExpr(expr ast.Expr) error {
	switch t := expr.(type) {
	case *ast.CallExpr:
		if t.Fun == nil {
			return nil
		}
		for _, a := range t.Args {
			switch at := a.(type) {
			case *ast.CallExpr:
				if sel, ok := t.Fun.(*ast.SelectorExpr); ok {
					return v.evalSelector(at, sel)
				}

				if err := v.evalArgs(at); err != nil {
					return err
				}
			case *ast.CompositeLit:
				for _, e := range at.Elts {
					if err := v.evalExpr(e); err != nil {
						return err
					}
				}
			}
		}
		if ft, ok := t.Fun.(*ast.SelectorExpr); ok {
			return v.evalSelector(t, ft)
		}
	case *ast.KeyValueExpr:
		return v.evalExpr(t.Value)
	}
	return nil
}

func (v *Visitor) evalArgs(expr ast.Expr) error {
	switch at := expr.(type) {
	case *ast.CompositeLit:
		for _, e := range at.Elts {
			if err := v.evalExpr(e); err != nil {
				return err
			}
		}
	case *ast.CallExpr:
		if at.Fun == nil {
			return nil
		}
		switch st := at.Fun.(type) {
		case *ast.SelectorExpr:
			if err := v.evalSelector(at, st); err != nil {
				return err
			}
		case *ast.Ident:
			return v.evalIdent(st)
		}
		for _, a := range at.Args {
			if err := v.evalArgs(a); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visitor) evalSelector(expr *ast.CallExpr, sel *ast.SelectorExpr) error {
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil
	}
	if x.Name == "pkger" {
		switch sel.Sel.Name {
		case "Walk":
			if len(expr.Args) != 2 {
				return fmt.Errorf("`New` requires two arguments")
			}

			zz := func(e ast.Expr) (string, error) {
				switch at := e.(type) {
				case *ast.Ident:
					switch at.Obj.Kind {
					case ast.Var:
						if as, ok := at.Obj.Decl.(*ast.AssignStmt); ok {
							return v.fromVariable(as)
						}
					case ast.Con:
						if vs, ok := at.Obj.Decl.(*ast.ValueSpec); ok {
							return v.fromConstant(vs)
						}
					}
					return "", v.evalIdent(at)
				case *ast.BasicLit:
					return at.Value, nil
				case *ast.CallExpr:
					return "", v.evalExpr(at)
				}
				return "", fmt.Errorf("can't handle %T", e)
			}

			k1, err := zz(expr.Args[0])
			if err != nil {
				return err
			}
			if err := v.addPath(k1); err != nil {
				return err
			}

			return nil
		case "Open":
			for _, e := range expr.Args {
				switch at := e.(type) {
				case *ast.Ident:
					switch at.Obj.Kind {
					case ast.Var:
						if as, ok := at.Obj.Decl.(*ast.AssignStmt); ok {
							v.addVariable("", as)
						}
					case ast.Con:
						if vs, ok := at.Obj.Decl.(*ast.ValueSpec); ok {
							v.addConstant("", vs)
						}
					}
					return v.evalIdent(at)
				case *ast.BasicLit:
					return v.addPath(at.Value)
				case *ast.CallExpr:
					return v.evalExpr(at)
				}
			}
		}
	}

	return nil
}

func (v *Visitor) evalIdent(i *ast.Ident) error {
	if i.Obj == nil {
		return nil
	}
	if s, ok := i.Obj.Decl.(*ast.AssignStmt); ok {
		return v.evalStmt(s)
	}
	return nil
}

func (v *Visitor) fromVariable(as *ast.AssignStmt) (string, error) {
	if len(as.Rhs) == 1 {
		if bs, ok := as.Rhs[0].(*ast.BasicLit); ok {
			return bs.Value, nil
		}
	}
	return "", fmt.Errorf("unable to find value from variable %v", as)
}

func (v *Visitor) addVariable(bn string, as *ast.AssignStmt) error {
	bv, err := v.fromVariable(as)
	if err != nil {
		return nil
	}
	if len(bn) == 0 {
		bn = bv
	}
	return v.addPath(bn)
}

func (v *Visitor) fromConstant(vs *ast.ValueSpec) (string, error) {
	if len(vs.Values) == 1 {
		if bs, ok := vs.Values[0].(*ast.BasicLit); ok {
			return bs.Value, nil
		}
	}
	return "", fmt.Errorf("unable to find value from constant %v", vs)
}

func (v *Visitor) addConstant(bn string, vs *ast.ValueSpec) error {
	if len(vs.Values) == 1 {
		if bs, ok := vs.Values[0].(*ast.BasicLit); ok {
			bv := bs.Value
			if len(bn) == 0 {
				bn = bv
			}
			return v.addPath(bn)
		}
	}
	return nil
}
