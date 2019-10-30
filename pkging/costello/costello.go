package costello

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

type AllFn func(ref *Ref) (pkging.Pkger, error)

func All(t *testing.T, ref *Ref, fn AllFn) {
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

	pkg, err := fn(ref)
	r.NoError(err)

	for n, tt := range tests {
		t.Run(fmt.Sprintf("%T/%s", pkg, n), func(st *testing.T) {
			st.Parallel()
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
	r.Equal(a.ModTime().Format(time.RFC3339), b.ModTime().Format(time.RFC3339))
	r.Equal(a.Mode(), b.Mode())
	r.Equal(a.Name(), b.Name())
	r.Equal(a.Size(), b.Size())
}
