package costello

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func WalkTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	r.NoError(LoadRef(ref, pkg))

	name := "public"

	fp := filepath.Join(ref.Dir, name)

	var exp []os.FileInfo
	err := filepath.Walk(fp, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		exp = append(exp, info)
		return nil
	})
	r.NoError(err)

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
		cmpFileInfo(t, info, act[i])
	}
}
