package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"sort"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
)

var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

func Parse(her here.Info) ([]here.Path, error) {

	src, err := fromSource(her)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func fromSource(her here.Info) ([]here.Path, error) {
	root := her.Dir
	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", root)
	}

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, root, nil, 0)
	if err != nil {
		return nil, err
	}
	pm := map[string]here.Path{}
	for _, pkg := range pkgs {
		for _, pf := range pkg.Files {
			f := &file{
				info:     her,
				fset:     fset,
				astFile:  pf,
				filename: pf.Name.Name,
				paths:    map[string]here.Path{},
			}
			f.decls = make(map[string]string)

			x, err := f.find()
			if err != nil {
				return nil, err
			}
			for i, pt := range x {
				if pt.Pkg == "/" || pt.Pkg == "" {
					pt.Pkg = her.ImportPath
					x[i] = pt
				}
				pt = x[i]
				pm[pt.String()] = pt
				err = pkger.Walk(pt.String(), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}
					p, err := her.Parse(path)
					if err != nil {
						return err
					}

					if pt.Name != "/" {
						p.Name = pt.Name + p.Name
					}
					pm[p.String()] = p
					return nil
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}
	var paths []here.Path

	for _, v := range pm {
		paths = append(paths, v)
	}

	sort.Slice(paths, func(i, j int) bool {
		return paths[i].String() < paths[j].String()
	})

	// for _, i := range her.Imports {
	// 	fmt.Println(">>>TODO parser/parser.go:237: i ", i)
	// }

	return paths, nil
}

// func importName(pkg *ast.File) (string, error) {
// 	var v visitor
// 	var name string
// 	var err error
// 	v = func(node ast.Node) ast.Visitor {
// 		if node == nil {
// 			return v
// 		}
// 		switch t := node.(type) {
// 		case *ast.ImportSpec:
// 			s, err := strconv.Unquote(t.Path.Value)
// 			if err != nil {
// 				err = err
// 				return nil
// 			}
// 			if s != "github.com/markbates/pkger" {
// 				if t.Name == nil {
// 					name = "pkger"
// 					return v
// 				}
// 			}
// 		default:
// 			// fmt.Printf("%#v\n", node)
// 		}
// 		return v
// 	}
// 	ast.Walk(v, pkg)
//
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if len(name) == 0 {
// 		return "", io.EOF
// 	}
// 	return name, nil
// }

type Results struct {
	Paths []here.Path
	Path  here.Path
}
