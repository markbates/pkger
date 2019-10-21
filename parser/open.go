package parser

import (
	"go/token"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
)

var _ Decl = OpenDecl{}

type OpenDecl struct {
	file  *File
	pos   token.Pos
	value string
}

func (d OpenDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d OpenDecl) Pos() (token.Pos, error) {
	if d.pos <= 0 {
		return -1, os.ErrNotExist
	}
	return d.pos, nil
}

func (d OpenDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d OpenDecl) Files() ([]*File, error) {

	pt, err := pkger.Parse(d.value)
	if err != nil {
		return nil, err
	}

	her, err := here.Package(pt.Pkg)
	if err != nil {
		return nil, err
	}

	fp := filepath.Join(her.Dir, pt.Name)

	osf, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}

	if osf.IsDir() {
		wd := WalkDecl{
			file:  d.file,
			pos:   d.pos,
			value: d.value,
		}
		return wd.Files()
	}

	var files []*File
	files = append(files, &File{
		Abs:  filepath.Join(her.Dir, pt.Name),
		Path: pt,
		Here: her,
	})

	return files, nil
}
