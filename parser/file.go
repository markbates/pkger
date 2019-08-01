package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"

	"github.com/markbates/errx"
)

type ParsedFile struct {
	File    string
	FileSet *token.FileSet
	Ast     *ast.File
}

// ParseFileMode ...
func ParseFileMode(f string, mode parser.Mode) (ParsedFile, error) {
	pf := ParsedFile{
		File:    f,
		FileSet: token.NewFileSet(),
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return pf, err
	}
	src := string(b)

	pff, err := parser.ParseFile(pf.FileSet, f, src, mode)
	if err != nil && errx.Unwrap(err) != io.EOF {
		return pf, err
	}
	pf.Ast = pff

	return pf, nil
}

// ParseFile ...
func ParseFile(f string) (ParsedFile, error) {
	return ParseFileMode(f, 0)
}
