package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/markbates/pkger/here"
)

// var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

func Parse(her here.Info) (Decls, error) {
	src, err := fromSource(her)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func fromSource(her here.Info) (Decls, error) {
	root := her.Dir
	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", root)
	}

	var decls Decls

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fset := token.NewFileSet()

		if !info.IsDir() {
			return nil
		}
		pkgs, err := parser.ParseDir(fset, path, nil, 0)
		if err != nil {
			return err
		}

		for _, pkg := range pkgs {
			for name, pf := range pkg.Files {
				f := &file{
					fset:     fset,
					astFile:  pf,
					filename: name,
					current:  her,
				}

				x, err := f.find()
				if err != nil {
					return err
				}
				decls = append(decls, x...)
			}
		}
		return nil
	})

	return decls, err
}
