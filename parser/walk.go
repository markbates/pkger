package parser

import (
	"encoding/json"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/here"
)

var _ Decl = WalkDecl{}

type WalkDecl struct {
	file  *File
	pos   token.Position
	value string
}

func (d WalkDecl) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d WalkDecl) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "pkger.Walk",
		"file":  d.file,
		"pos":   d.pos,
		"value": d.value,
	})
}

func (d WalkDecl) File() (*File, error) {
	if d.file == nil {
		return nil, os.ErrNotExist
	}
	return d.file, nil
}

func (d WalkDecl) Position() (token.Position, error) {
	return d.pos, nil
}

func (d WalkDecl) Value() (string, error) {
	if d.value == "" {
		return "", os.ErrNotExist
	}
	return d.value, nil
}

func (d WalkDecl) Files(virtual map[string]string) ([]*File, error) {
	var files []*File

	her := d.file.Here
	pt := d.file.Path

	root := filepath.Join(her.Module.Dir, pt.Name)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			her, err = here.Dir(path)
			if err != nil {
				return err
			}
		}

		n := strings.TrimPrefix(path, her.Module.Dir)

		pt, err := her.Parse(n)
		if err != nil {
			return err
		}

		if _, ok := virtual[n]; ok {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, &File{
			Abs:  path,
			Path: pt,
			Here: her,
		})
		return nil
	})

	return files, err
}
