package parser

import (
	"encoding/json"
	"fmt"
	"go/token"
	"os"
)

var _ Decl = IncludeDecl{}

type IncludeDecl struct {
	file  *File
	pos   token.Position
	value string
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
