package mem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Add(t *testing.T) {
	r := require.New(t)

	cur, err := here.Package("github.com/markbates/pkger")
	r.NoError(err)

	pkg, err := mem.New(cur)
	r.NoError(err)

	var exp []os.FileInfo
	root := filepath.Join(cur.Dir, "pkging", "pkgtest", "testdata", "ref")
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		exp = append(exp, info)
		return pkg.Add(f)
	})
	r.NoError(err)

	var act []os.FileInfo
	err = pkg.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, info)
		return nil
	})
	r.NoError(err)

	r.Len(act, len(exp))

	for i, e := range exp {
		a := act[i]

		r.Equal(e.Name(), a.Name())
		r.Equal(e.Size(), a.Size())
		r.Equal(e.Mode(), a.Mode())
		r.Equal(e.IsDir(), a.IsDir())
	}

}
