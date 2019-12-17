package parser

import (
	"fmt"
	"go/token"
	"sort"
)

type Decl interface {
	File() (*File, error)
	Position() (token.Position, error)
	Value() (string, error)
}

type Filer interface {
	Files(map[string]string) ([]*File, error)
}

type Virtualer interface {
	VirtualPaths() []string
}

type Decls []Decl

func (decls Decls) Files() ([]*File, error) {
	m := map[string]*File{}
	v := map[string]string{}

	for _, d := range decls {
		if vt, ok := d.(Virtualer); ok {
			for _, s := range vt.VirtualPaths() {
				v[s] = s
			}
		}

		fl, ok := d.(Filer)
		if !ok {
			continue
		}

		files, err := fl.Files(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, d)
		}

		for _, f := range files {
			m[f.Abs] = f
			v[f.Abs] = f.Abs
		}
	}

	var files []*File
	for _, f := range m {
		files = append(files, f)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path.String() < files[j].Path.String()
	})
	return files, nil
}
