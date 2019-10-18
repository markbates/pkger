package mem_test

import (
	"bytes"
	"os"
	"sort"
	"testing"

	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/markbates/pkger/pkging/stuffing"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Embedding(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	paths, err := parser.Parse(app.Info)
	r.NoError(err)

	ps := make([]string, len(paths))
	for i, p := range paths {
		ps[i] = p.String()
	}

	r.Equal(app.Paths.Parser, ps)

	base, err := mem.New(app.Info)
	r.NoError(err)

	disk, err := stdos.New(app.Info)
	r.NoError(err)

	err = disk.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := disk.Open(path)
		if err != nil {
			return err
		}
		return base.Add(f)

	})
	r.NoError(err)

	var res []string
	err = base.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)
	r.Equal(app.Paths.Root, res)

	res = []string{}
	err = base.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)
	r.Equal(app.Paths.Public, res)

	bb := &bytes.Buffer{}

	err = stuffing.Stuff(bb, app.Info, paths)
	r.NoError(err)

	pkg, err := mem.UnmarshalEmbed(bb.Bytes())
	r.NoError(err)

	res = []string{}
	err = pkg.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)

	exp := append(app.Paths.Public, "app:/")
	sort.Strings(exp)
	r.Equal(exp, res)

	res = []string{}
	err = pkg.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)
	r.Equal(app.Paths.Public, res)
}
