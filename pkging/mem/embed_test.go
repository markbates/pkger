package mem_test

import (
	"bytes"
	"os"
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

	info, err := pkgtest.App()
	r.NoError(err)

	paths, err := parser.Parse(info)
	r.NoError(err)

	ps := make([]string, len(paths))
	for i, p := range paths {
		ps[i] = p.String()
	}

	r.Equal(pkgtest.AppPaths, ps)

	base, err := mem.New(info)
	r.NoError(err)

	disk, err := stdos.New(info)
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
	r.Equal(rootWalk, res)

	res = []string{}
	err = base.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)
	r.Equal(pubWalk[1:], res)

	bb := &bytes.Buffer{}

	err = stuffing.Stuff(bb, info, paths)
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
	r.Equal(rootWalk, res)

	res = []string{}
	err = pkg.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)
	r.Equal(pubWalk, res)
}

var pubWalk = []string{
	"app:/",
	"app:/public",
	"app:/public/images",
	"app:/public/images/img1.png",
	"app:/public/images/img2.png",
	"app:/public/index.html",
}

var rootWalk = []string{
	"app:/",
	"app:/go.mod",
	"app:/main.go",
	"app:/public",
	"app:/public/images",
	"app:/public/images/img1.png",
	"app:/public/images/img2.png",
	"app:/public/index.html",
	"app:/templates",
	"app:/templates/a.txt",
	"app:/templates/b",
	"app:/templates/b/b.txt",
}
