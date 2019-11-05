package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/stretchr/testify/require"
)

func Test_Parser_Ref(t *testing.T) {
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	disk, err := stdos.New(ref.Info)
	r.NoError(err)

	_, err = pkgtest.LoadFiles("/", ref, disk)
	r.NoError(err)

	res, err := parser.Parse(ref.Info)

	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)
	r.Len(files, 23)

	for _, f := range files {
		if f.Path.Pkg == ref.Module.Path {
			r.True(strings.HasPrefix(f.Abs, ref.Dir), "%q %q", f.Abs, ref.Dir)
		} else {
			r.False(strings.HasPrefix(f.Abs, ref.Dir), "%q %q", f.Abs, ref.Dir)
		}
	}
}

func Test_Parser_Example_HTTP(t *testing.T) {
	r := require.New(t)

	cur, err := here.Current()
	r.NoError(err)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	root := filepath.Join(cur.Dir, "examples", "http", "pkger")
	r.NoError(os.Chdir(root))

	her, err := here.Dir(".")
	r.NoError(err)

	res, err := parser.Parse(her)
	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)
	r.Len(files, 5)

	for _, f := range files {
		r.True(strings.HasPrefix(f.Abs, her.Dir), "%q %q", f.Abs, her.Dir)
		r.True(strings.HasPrefix(f.Path.Name, "/public"), "%q %q", f.Path.Name, "/public")
	}
}
