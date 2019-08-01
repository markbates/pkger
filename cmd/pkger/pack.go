package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"text/template"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/paths"
	"github.com/markbates/pkger/pkgs"
)

const outName = "pkged.go"

func pack(args []string) error {
	info, err := pkgs.Current()
	if err != nil {
		return err
	}

	fp := info.FilePath(outName)
	os.RemoveAll(fp)
	if len(args) == 0 {
		args = append(args, "")
	}

	res, err := parser.Parse(args[0])
	if err != nil {
		return err
	}

	if err := Package(fp, res.Paths); err != nil {
		return err
	}

	return nil
}

func Package(p string, paths []paths.Path) error {
	os.RemoveAll(p)

	var files []*pkger.File
	for _, p := range paths {
		f, err := pkger.Open(p.String())
		if err != nil {
			return err
		}
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			continue
		}
		f.Close()
		files = append(files, f)
	}

	bb := &bytes.Buffer{}
	gz := gzip.NewWriter(bb)
	defer gz.Close()

	enc := json.NewEncoder(gz)
	if err := enc.Encode(files); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	s := base64.StdEncoding.EncodeToString(bb.Bytes())
	//
	d := struct {
		Pkg  string
		Data string
	}{
		Pkg:  "jeremy",
		Data: strconv.Quote(s),
	}

	f, err := os.Create(p)
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

	// fmt.Fprintf(f, "package jeremy\n\n")
	// fmt.Fprintf(f, "var _ = func() error {\n")
	// fmt.Fprintf(f, "const data = %q\n", s)
	// fmt.Fprintf(f, "return nil\n")
	// fmt.Fprintf(f, "}\n")

	// fmt.Println(bb.String())
	// f, err := os.Create(".pkger.index.json")
	// if err != nil {
	// 	return err
	// }
	//
	// enc := json.NewEncoder(f)
	// if err := enc.Encode(files); err != nil {
	// 	return err
	// }
	return nil
}

const outTmpl = `package {{.Pkg}}

var _ = func() error {
		const data = {{.Data}}
		return nil
}
`
