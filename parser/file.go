package parser

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"

	"github.com/markbates/pkger/here"
)

type File struct {
	Abs  string // full path on disk to file
	Path here.Path
	Here here.Info
}

func (f File) String() string {
	b, _ := json.MarshalIndent(f, "", " ")
	return string(b)
}

type parsedFile struct {
	File    string
	FileSet *token.FileSet
	Ast     *ast.File
}

// parseFileMode ...
func parseFileMode(f string, mode parser.Mode) (parsedFile, error) {
	pf := parsedFile{
		File:    f,
		FileSet: token.NewFileSet(),
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return pf, err
	}
	src := string(b)

	pff, err := parser.ParseFile(pf.FileSet, f, src, mode)
	if err != nil && err != io.EOF {
		return pf, err
	}
	pf.Ast = pff

	return pf, nil
}

// parseFile ...
func parseFile(f string) (parsedFile, error) {
	return parseFileMode(f, 0)
}
