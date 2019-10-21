package mem_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Embedding(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	res, err := parser.Parse(app.Info)
	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)

	ps := make([]string, len(files))
	for i, f := range files {
		ps[i] = f.Path.String()
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

	var act []string
	err = base.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, path)
		return nil
	})

	r.NoError(err)
	r.Equal(app.Paths.Root, act)

	act = []string{}
	err = base.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, path)
		return nil
	})

	r.NoError(err)
	r.Equal(app.Paths.Public, act)

	bb := &bytes.Buffer{}

	err = pkgutil.Stuff(bb, app.Info, res)
	r.NoError(err)

	pkg := &mem.Pkger{}
	err = pkg.UnmarshalEmbed(bb.Bytes())
	r.NoError(err)

	act = []string{}
	err = pkg.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, path)
		return nil
	})

	r.NoError(err)

	r.Equal(app.Paths.Public, act)

	act = []string{}
	err = pkg.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, path)
		return nil
	})

	r.NoError(err)
	r.Equal(app.Paths.Public, act)
}
