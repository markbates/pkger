package here_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Dir(t *testing.T) {
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)

	root := ref.Dir

	r.NoError(err)
	defer os.RemoveAll(root)

	public := filepath.Join(root, "public")
	r.NoError(os.MkdirAll(public, 0755))

	gf := filepath.Join(root, "cmd", "main.go")
	r.NoError(os.MkdirAll(filepath.Dir(gf), 0755))

	f, err := os.Create(gf)
	r.NoError(err)

	_, err = f.Write([]byte("package main"))
	r.NoError(err)

	r.NoError(f.Close())

	table := []struct {
		in  string
		err bool
	}{
		{in: root, err: false},
		{in: public, err: false},
		{in: gf, err: false},
		{in: filepath.Join(root, "."), err: false},
		{in: "/unknown", err: true},
	}
	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			here.ClearCache()
			r := require.New(st)

			info, err := here.Dir(tt.in)
			if tt.err {
				r.Error(err)
				return
			}
			r.NoError(err)

			r.NotZero(info)
			r.NotZero(info.Dir)
			r.NotZero(info.Name)
			r.NotZero(info.Module)

		})
	}
}
