package mem_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Add(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	pkg, err := mem.New(app.Info)
	r.NoError(err)

	root := app.Info.Dir

	var exp []os.FileInfo
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		exp = append(exp, info)

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		return pkg.Add(f)

	})

	r.NoError(err)

	sort.Slice(exp, func(i, j int) bool {
		return exp[i].Name() < exp[j].Name()
	})

	var act []os.FileInfo
	err = pkg.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, info)
		return nil
	})
	r.NoError(err)

	sort.Slice(act, func(i, j int) bool {
		return act[i].Name() < act[j].Name()
	})

	r.Len(act, len(exp))

	for i, e := range exp {
		a := act[i]

		r.Equal(e.Name(), a.Name())
		r.Equal(e.Size(), a.Size())
		r.Equal(e.Mode(), a.Mode())
		r.Equal(e.IsDir(), a.IsDir())
	}

}
