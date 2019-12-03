package pkgutil

import (
	"bytes"
	"os"
	"testing"

	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/stretchr/testify/require"
)

func Test_Stuff(t *testing.T) {
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(ref.Dir)

	disk, err := stdos.New(ref.Info)
	r.NoError(err)

	infos, err := pkgtest.LoadFiles("/", ref, disk)
	r.NoError(err)
	r.Len(infos, 34)

	decls, err := parser.Parse(ref.Info)
	r.NoError(err)

	r.Len(decls, 11)

	files, err := decls.Files()
	r.NoError(err)

	for _, f := range files {
		if f.Path.Pkg == ref.Module.Path {
			r.Equal("app", f.Path.Pkg)
		} else {
			r.NotEqual("app", f.Path.Pkg)
		}
	}

	r.Len(files, 25)

	bb := &bytes.Buffer{}

	err = Stuff(bb, ref.Info, decls)
	r.NoError(err)

	pkg, err := mem.UnmarshalEmbed(bb.Bytes())
	r.NoError(err)

	pkgtest.CurrentTest(t, ref, pkg)
	pkgtest.InfoTest(t, ref, pkg)
	pkgtest.OpenTest(t, ref, pkg)
	pkgtest.WalkTest(t, ref, pkg)

	_, err = pkg.Stat("/go.mod")
	r.NoError(err)

	_, err = pkg.Stat("/public/index.html")
	r.NoError(err)
}
