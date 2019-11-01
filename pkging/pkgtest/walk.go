package pkgtest

import (
	"os"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func WalkTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	exp, err := LoadFiles("/public", ref, pkg)
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	name := "public"

	var act []os.FileInfo
	err = pkg.Walk("/"+name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, info)
		return nil
	})

	r.NoError(err)

	r.Len(act, len(exp))

	for i, info := range exp {
		CmpFileInfo(t, info, act[i])
	}
}
