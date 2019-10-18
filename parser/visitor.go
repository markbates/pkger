package parser

import (
	"go/ast"
	"go/token"
	"sort"
	"strconv"
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
	decls    map[string]string
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.astFile)
}

func (f *file) find() ([]string, error) {
	// if err := f.findDecals(); err != nil {
	// 	return nil, err
	// }
	if err := f.findOpenCalls(); err != nil {
		return nil, err
	}

	if err := f.findWalkCalls(); err != nil {
		return nil, err
	}

	if err := f.findImportCalls(); err != nil {
		return nil, err
	}

	var paths []string
	for _, p := range f.decls {
		paths = append(paths, p)
	}
	sort.Slice(paths, func(a, b int) bool {
		return paths[a] < paths[b]
	})
	return paths, nil
}

func (f *file) findDecals() error {
	// iterate over all declarations
	for _, d := range f.astFile.Decls {

		// log.Printf("#%d Decl: %+v\n", i, d)

		// only interested in generic declarations
		if genDecl, ok := d.(*ast.GenDecl); ok {

			// handle const's and vars
			if genDecl.Tok == token.CONST || genDecl.Tok == token.VAR {

				// there may be multiple
				// i.e. const ( ... )
				for _, cDecl := range genDecl.Specs {

					// havn't find another kind of spec then value but better check
					if vSpec, ok := cDecl.(*ast.ValueSpec); ok {
						// log.Printf("const ValueSpec: %+v\n", vSpec)

						// iterate over Name/Value pair
						for i := 0; i < len(vSpec.Names); i++ {
							// TODO: only basic literals work currently
							if i > len(vSpec.Values) || len(vSpec.Values) == 0 {
								break
							}
							switch v := vSpec.Values[i].(type) {
							case *ast.BasicLit:

								if e := f.addNode(v); e != nil {
									return e
								}
								// f.decls[vSpec.Names[i].Name] = v.Value
							default:
								// log.Printf("Name: %s - Unsupported ValueSpec: %+v\n", vSpec.Names[i].Name, v)
							}
						}
					}
				}
			}

		}
	}

	return nil
}

func (f *file) addNode(node ast.Node) error {
	switch x := node.(type) {
	case *ast.BasicLit:
		return f.add(x.Value)
	case *ast.Ident:
		return f.add(x.Name)
	default:
	}
	return nil
}

func (f *file) add(s string) error {
	s, err := strconv.Unquote(s)
	if err != nil {
		return err
	}
	if _, ok := f.decls[s]; !ok {
		// fmt.Println(">>>TODO parser/visitor.go:98: s ", s)
		f.decls[s] = s
	}
	return nil
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

		// fmt.Println(">>>TODO parser/visitor.go:138: findOpenCalls ", ce.Args[0])
		if e := f.addNode(ce.Args[0]); e != nil {
			err = e
			return false
		}

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

		// fmt.Println(">>>TODO parser/visitor.go:138: findWalkCalls ", ce.Args[0])
		if e := f.addNode(ce.Args[0]); e != nil {
			err = e
			return false
		}

		return true
	})
	return err
}

func (f *file) findImportCalls() error {
	var err error
	f.walk(func(node ast.Node) bool {
		// ce, ok := node.(*ast.ImportSpec)
		// if !ok {
		// 	return true
		// }

		// s, err := strconv.Unquote(ce.Path.Value)
		// if err != nil {
		// 	return false
		// }
		// fmt.Println(">>>TODO parser/visitor.go:215: s ", s)
		// info, err := here.Package(s)
		// if err != nil {
		// 	return false
		// }
		// fmt.Println(">>>TODO parser/visitor.go:216: info ", info)
		// res, err := Parse(info)
		// if err != nil {
		// 	return false
		// }
		// fmt.Println(">>>TODO parser/visitor.go:224: res ", res)
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
