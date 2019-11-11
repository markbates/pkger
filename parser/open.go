package parser

import (
	"encoding/json"
	"go/token"
	"os"
	"path/filepath"
)

var _ Decl = OpenDecl{}

type OpenDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d OpenDecl) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d OpenDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.Open",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d OpenDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d OpenDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d OpenDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d OpenDecl) Files(virtual map[string]string) ([]*File, error) {
	if _, ok := virtual[d.value]; ok {
		return nil, nil
	}

	her := d.file.Here
	pt := d.file.Path

	fp := filepath.Join(her.Module.Dir, pt.Name)

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
		return wd.Files(virtual)
	}

	var files []*File
	files = append(files, d.file)

	return files, nil
}
