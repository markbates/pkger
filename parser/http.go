package parser

import (
	"encoding/json"
	"go/token"
	"os"
)

var _ Decl = HTTPDecl{}

type HTTPDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d HTTPDecl) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d HTTPDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.Dir",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d HTTPDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d HTTPDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d HTTPDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d HTTPDecl) Files(virtual map[string]string) ([]*File, error) {
	od := OpenDecl{
		file:  d.file,
		pos:   d.pos,
		value: d.value,
	}

	return od.Files(virtual)
}
