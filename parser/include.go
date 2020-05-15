package parser

import (
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"path/filepath"

	"github.com/markbates/pkger/here"
)

var _ Decl = IncludeDecl{}

type IncludeDecl struct {
	file  *File
	pos   token.Position
	value string
}

func NewInclude(her here.Info, inc string) (IncludeDecl, error) {
	var id IncludeDecl
	pt, err := her.Parse(inc)
	if err != nil {
		return id, err
	}

	if pt.Pkg != her.ImportPath {
		her, err = here.Package(pt.Pkg)
		if err != nil {
			return id, err
		}
	}

	abs := filepath.Join(her.Module.Dir, pt.Name)

	f := &File{
		Abs:  abs,
		Path: pt,
		Here: her,
	}

	return IncludeDecl{
		value: inc,
		file:  f,
	}, nil
}

func (d IncludeDecl) String() string {
	return fmt.Sprintf("pkger.Include(%q)", d.value)
}

func (d IncludeDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.Include",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d IncludeDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d IncludeDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d IncludeDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d IncludeDecl) Files(virtual map[string]string) ([]*File, error) {
	od := OpenDecl{
		file:  d.file,
		pos:   d.pos,
		value: d.value,
	}

	return od.Files(virtual)
}
