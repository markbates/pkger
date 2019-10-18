package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/stdos"
)

// var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

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
				fset:     fset,
				astFile:  pf,
				filename: pf.Name.Name,
				decls:    map[string]string{},
			}

			x, err := f.find()
			if err != nil {
				return nil, err
			}
			for _, dl := range x {
				pt, err := her.Parse(dl)
				if err != nil {
					return nil, err
				}
				res, err := fromPath(pt)
				if err != nil {
					return nil, err
				}
				for _, p := range res {
					pm[p.String()] = p
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

	return paths, nil
}

func fromPath(pt here.Path) ([]here.Path, error) {
	var paths []here.Path

	her, err := here.Package(pt.Pkg)
	if err != nil {
		return nil, err
	}

	pkg, err := stdos.New(her)
	if err != nil {
		return nil, err
	}

	root := here.Path{
		Pkg:  pt.Pkg,
		Name: strings.Replace(filepath.Dir(pt.Name), "\\", "/", -1),
	}
	paths = append(paths, root)
	err = pkg.Walk(pt.Name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		p, err := her.Parse(path)
		if err != nil {
			return err
		}
		paths = append(paths, p)

		return nil
	})

	return paths, nil
}
