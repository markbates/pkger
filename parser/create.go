package parser

import (
	"encoding/json"
	"go/token"
	"os"
)

var _ Decl = CreateDecl{}

type CreateDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d CreateDecl) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d CreateDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.Create",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d CreateDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d CreateDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d CreateDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d CreateDecl) VirtualPaths() []string {
	return []string{d.value}
}
