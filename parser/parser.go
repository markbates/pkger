package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/markbates/oncer"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

func Parse(her here.Info) (Results, error) {
	var r Results
	var err error

	oncer.Do(her.ImportPath, func() {
		pwd, err := os.Getwd()
		if err != nil {
			return
		}
		defer os.Chdir(pwd)

		fmt.Println("cd: ", her.Dir, her.ImportPath)
		os.Chdir(her.Dir)

		// 1. search for .go files in/imported by `her.ImportPath`
		src, err := fromSource(her)
		if err != nil {
			return
		}
		fmt.Println(">>>TODO parser/parser.go:30: src ", src)

		// 2. parse .go ast's for `pkger.*` calls
		// 3. find path's in those files
		// 4. walk folders in those paths and add to results
	})

	return r, err
}

func fromSource(her here.Info) ([]pkging.Path, error) {
	fmt.Println(">>>TODO parser/parser.go:201: her.ImportPath ", her.ImportPath)
	fmt.Println(her)
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

	var paths []pkging.Path
	for _, pkg := range pkgs {
		for _, pf := range pkg.Files {
			f := &file{fset: fset, astFile: pf, filename: pf.Name.Name}
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
				paths = append(paths, x[i])
			}
		}
	}

	for _, i := range her.Imports {
		fmt.Println(">>>TODO parser/parser.go:237: i ", i)
	}

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
	Paths []pkging.Path
	Path  pkging.Path
}
