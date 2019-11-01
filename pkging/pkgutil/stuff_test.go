package pkgutil

import (
	"bytes"
	"os"
	"testing"

	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/costello"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/stretchr/testify/require"
)

func Test_Stuff(t *testing.T) {
	r := require.New(t)

	ref, err := costello.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	disk, err := stdos.New(ref.Info)
	r.NoError(err)

	infos, err := costello.LoadFiles("/", ref, disk)
	r.NoError(err)
	r.Len(infos, 34)

	decls, err := parser.Parse(ref.Info)
	r.NoError(err)

	r.Len(decls, 9)

	files, err := decls.Files()
	r.NoError(err)

	for _, f := range files {
		r.Equal("app", f.Path.Pkg)
	}

	r.Len(files, 22)

	bb := &bytes.Buffer{}

	err = Stuff(bb, ref.Info, decls)
	r.NoError(err)

	pkg, err := mem.UnmarshalEmbed(bb.Bytes())
	r.NoError(err)

	costello.CurrentTest(t, ref, pkg)
	costello.InfoTest(t, ref, pkg)
	costello.OpenTest(t, ref, pkg)
	costello.WalkTest(t, ref, pkg)

	_, err = pkg.Stat("/go.mod")
	r.NoError(err)

	_, err = pkg.Stat("/public/index.html")
	r.NoError(err)
}
