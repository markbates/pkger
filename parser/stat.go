package parser

import (
	"encoding/json"
	"fmt"
	"go/token"
	"os"
)

var _ Decl = StatDecl{}

type StatDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d StatDecl) String() string {
	return fmt.Sprintf("pkger.Stat(%q)", d.value)
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
	od := OpenDecl{
		file:  d.file,
		pos:   d.pos,
		value: d.value,
	}

	return od.Files(virtual)
}
