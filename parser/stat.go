package parser

import (
	"encoding/json"
	"go/token"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
)

var _ Decl = StatDecl{}

type StatDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d StatDecl) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d StatDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.Stat",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d StatDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d StatDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d StatDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d StatDecl) Files(virtual map[string]string) ([]*File, error) {
	if _, ok := virtual[d.value]; ok {
		return nil, nil
	}

	pt, err := pkger.Parse(d.value)
	if err != nil {
		return nil, err
	}

	her, err := here.Package(pt.Pkg)
	if err != nil {
		return nil, err
	}

	var files []*File
	files = append(files, &File{
		Abs:  filepath.Join(her.Module.Dir, pt.Name),
		Path: pt,
		Here: her,
	})

	return files, nil
}
