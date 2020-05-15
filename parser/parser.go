package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/markbates/pkger/here"
)

var defaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "testdata"}

func New(her here.Info) (*Parser, error) {
	return &Parser{
		Info:  her,
		decls: map[string]Decls{},
	}, nil
}

type Parser struct {
	here.Info
	decls    map[string]Decls
	once     sync.Once
	includes []string
	err      error
}

func Parse(her here.Info, includes ...string) (Decls, error) {
	p, err := New(her)
	if err != nil {
		return nil, err
	}
	p.includes = includes

	return p.Decls()
}

func (p *Parser) ParseSource(source Source, mode parser.Mode) (*ParsedSource, error) {
	pf := &ParsedSource{
		Source:  source,
		FileSet: token.NewFileSet(),
	}

	b, err := ioutil.ReadFile(source.Abs)
	if err != nil {
		return nil, err
	}
	src := string(b)

	pff, err := parser.ParseFile(pf.FileSet, source.Abs, src, mode)
	if err != nil && err != io.EOF {
		return nil, err
	}
	pf.Ast = pff

	return pf, nil
}

func (p *Parser) ParseFile(abs string, mode parser.Mode) (*ParsedSource, error) {
	s := Source{
		Abs: abs,
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", abs)
	}

	dir := filepath.Dir(abs)

	s.Here, err = here.Dir(dir)
	if err != nil {
		return nil, err
	}

	s.Path, err = s.Here.Parse(strings.TrimPrefix(abs, dir))

	return p.ParseSource(s, 0)
}

func (p *Parser) ParseDir(abs string, mode parser.Mode) ([]*ParsedSource, error) {
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", abs)
	}

	her, err := here.Dir(abs)
	if err != nil {
		return nil, fmt.Errorf("%s: here.Dir failed %s", err, abs)
	}

	pt, err := her.Parse(strings.TrimPrefix(abs, filepath.Dir(abs)))
	if err != nil {
		return nil, fmt.Errorf("%s: here.Parse failed %s", err, abs)
	}

	filter := func(f os.FileInfo) bool {
		return !f.IsDir()
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, abs, filter, 0)
	if err != nil {
		return nil, fmt.Errorf("%s: ParseDir failed %s", err, abs)
	}

	var srcs []*ParsedSource
	for _, pkg := range pkgs {
		for name, pf := range pkg.Files {
			s := &ParsedSource{
				Source: Source{
					Abs:  name,
					Path: pt,
					Here: her,
				},
				FileSet: fset,
				Ast:     pf,
			}
			srcs = append(srcs, s)
		}
	}

	return srcs, nil
}

func (p *Parser) Decls() (Decls, error) {
	if err := p.parse(); err != nil {
		return nil, err
	}

	var decls Decls
	orderedNames := []string{
		"MkdirAll",
		"Create",
		"Include",
		"Stat",
		"Open",
		"Dir",
		"Walk",
	}

	for _, n := range orderedNames {
		decls = append(decls, p.decls[n]...)
	}

	return decls, nil
}

func (p *Parser) DeclsMap() (map[string]Decls, error) {
	err := p.Parse()
	return p.decls, err
}

func (p *Parser) Parse() error {
	(&p.once).Do(func() {
		p.err = p.parse()
	})
	return p.err
}

func (p *Parser) parse() error {
	p.decls = map[string]Decls{}

	root := p.Dir

	if err := p.parseIncludes(); err != nil {
		return err
	}

	fi, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory", root)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		for _, x := range defaultIgnoredFolders {
			if strings.HasPrefix(base, x) {
				return filepath.SkipDir
			}
		}

		srcs, err := p.ParseDir(path, 0)
		if err != nil {
			return fmt.Errorf("%s: %s", err, path)
		}

		for _, src := range srcs {
			mm, err := src.DeclsMap()
			if err != nil {
				return fmt.Errorf("%s: %s", err, src.Abs)
			}
			for k, v := range mm {
				p.decls[k] = append(p.decls[k], v...)
			}
		}

		return nil
	})

	return err
}

func (p *Parser) parseIncludes() error {
	for _, i := range p.includes {
		d, err := NewInclude(p.Info, i)
		if err != nil {
			return err
		}

		p.decls["Include"] = append(p.decls["Include"], d)
	}
	return nil
}
