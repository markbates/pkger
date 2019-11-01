package pkgtest

import (
	"fmt"
	"os"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

type AllFn func(ref *Ref) (pkging.Pkger, error)

func All(t *testing.T, fn AllFn) {
	r := require.New(t)

	type tf func(*testing.T, *Ref, pkging.Pkger)

	tests := map[string]tf{
		"Create":    CreateTest,
		"Current":   CurrentTest,
		"HTTP":      HTTPTest,
		"Info":      InfoTest,
		"MkdirAll":  MkdirAllTest,
		"Open":      OpenTest,
		"Remove":    RemoveTest,
		"RemoveAll": RemoveAllTest,
		"Stat":      StatTest,
		"Walk":      WalkTest,
	}

	ref, err := NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	pkg, err := fn(ref)
	r.NoError(err)

	for n, tt := range tests {
		t.Run(fmt.Sprintf("%T/%s", pkg, n), func(st *testing.T) {
			st.Parallel()

			r := require.New(st)

			ref, err := NewRef()
			r.NoError(err)
			defer os.RemoveAll(ref.Dir)

			pkg, err := fn(ref)
			r.NoError(err)

			tt(st, ref, pkg)
		})
	}

}

func CmpFileInfo(t *testing.T, a os.FileInfo, b os.FileInfo) {
	t.Helper()

	r := require.New(t)
	r.Equal(a.IsDir(), b.IsDir())
	r.Equal(a.Name(), b.Name())
	r.NotZero(b.ModTime())

	if a.IsDir() {
		r.True(b.Mode().IsDir(), b.Mode().String())
		return
	}

	r.True(b.Mode().IsRegular(), b.Mode().String())
}

func cmpHereInfo(t *testing.T, a here.Info, b here.Info) {
	t.Helper()

	r := require.New(t)

	r.NotZero(a)
	r.NotZero(b)

	r.Equal(a.ImportPath, b.ImportPath)
	r.Equal(a.Name, b.Name)

	am := a.Module
	bm := b.Module

	r.Equal(am.Path, bm.Path)
	r.Equal(am.Main, bm.Main)
	r.Equal(am.GoVersion, bm.GoVersion)
}
