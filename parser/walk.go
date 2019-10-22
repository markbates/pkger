package parser

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
)

var _ Decl = WalkDecl{}

type WalkDecl struct {
	file  *File
	pos   token.Pos
	value string
}

func (d WalkDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d WalkDecl) Pos() (token.Pos, error) {
	if d.pos <= 0 {
		return -1, os.ErrNotExist
	}
	return d.pos, nil
}

func (d WalkDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d WalkDecl) Files() ([]*File, error) {
	pt, err := pkger.Parse(d.value)
	if err != nil {
		return nil, err
	}

	cur, err := here.Package(pt.Pkg)
	if err != nil {
		return nil, err
	}

	root := filepath.Join(cur.Dir, pt.Name)

	var files []*File
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		n := strings.TrimPrefix(path, cur.Dir)

		pt, err := pkger.Parse(n)
		if err != nil {
			return err
		}
		pt.Pkg = cur.ImportPath

		files = append(files, &File{
			Abs:  path,
			Path: pt,
			Here: cur,
		})
		return nil
	})

	return files, err
}
