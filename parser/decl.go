package parser

import (
	"go/token"
	"sort"
)

type Decl interface {
	File() (*File, error)
	Pos() (token.Pos, error)
	Value() (string, error)
}

type Filer interface {
	Files() ([]*File, error)
}

type Decls []Decl

func (decls Decls) Files() ([]*File, error) {
	m := map[string]*File{}

	for _, d := range decls {
		fl, ok := d.(Filer)
		if !ok {
			continue
		}

		files, err := fl.Files()
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			m[f.Path.String()] = f
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
