package parser

import (
	"encoding/json"
	"go/token"
	"os"
)

var _ Decl = MkdirAllDecl{}

type MkdirAllDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d MkdirAllDecl) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d MkdirAllDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.MkdirAll",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d MkdirAllDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d MkdirAllDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d MkdirAllDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d MkdirAllDecl) VirtualPaths() []string {
	return []string{d.value}
}
