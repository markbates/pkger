package costello

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
		"Open":      OpenTest,
		"Stat":      StatTest,
		"Create":    CreateTest,
		"Current":   CurrentTest,
		"Info":      InfoTest,
		"MkdirAll":  MkdirAllTest,
		"Remove":    RemoveTest,
		"RemoveAll": RemoveAllTest,
		"Walk":      WalkTest,
	}

	ref, err := NewRef()
	r.NoError(err)

	pkg, err := fn(ref)
	r.NoError(err)

	for n, tt := range tests {
		t.Run(fmt.Sprintf("%T/%s", pkg, n), func(st *testing.T) {
			st.Parallel()

			ref, err := NewRef()
			r.NoError(err)

			pkg, err := fn(ref)
			r.NoError(err)

			tt(st, ref, pkg)
		})
	}

}

func cmpFileInfo(t *testing.T, a os.FileInfo, b os.FileInfo) {
	t.Helper()

	r := require.New(t)

	r.Equal(a.IsDir(), b.IsDir())
	r.Equal(a.Mode().String(), b.Mode().String())
	r.Equal(a.Name(), b.Name())
	r.Equal(a.Size(), b.Size())
	r.NotZero(b.ModTime())
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
