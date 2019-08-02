package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"os"
	"strconv"
	"text/template"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

const outName = "pkged.go"

func pack(args []string) error {
	info, err := pkger.Current()
	if err != nil {
		return err
	}

	fp := info.FilePath(outName)
	os.RemoveAll(fp)

	res, err := parser.Parse(info.Dir)
	if err != nil {
		return err
	}

	if err := Package(fp, res.Paths); err != nil {
		return err
	}

	return nil
}

func Package(out string, paths []pkger.Path) error {
	os.RemoveAll(out)

	bb := &bytes.Buffer{}
	gz := gzip.NewWriter(bb)
	defer gz.Close()

	if err := pkger.Pack(gz, paths); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	s := base64.StdEncoding.EncodeToString(bb.Bytes())

	c, err := pkger.Current()
	if err != nil {
		return err
	}
	d := struct {
		Pkg  string
		Data string
	}{
		Pkg:  c.Name,
		Data: strconv.Quote(s),
	}

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New(outName).Parse(outTmpl)
	if err != nil {
		return err
	}
	if err := t.Execute(f, d); err != nil {
		return err
	}
	return nil
}

const outTmpl = `package {{.Pkg}}

import "github.com/markbates/pkger"

var _ = pkger.Unpack({{.Data}})
`
