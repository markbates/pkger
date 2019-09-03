package stdos

import (
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Pkger(t *testing.T) {
	r := require.New(t)

	mypkging, err := New()
	r.NoError(err)

	mypkging.current.Dir = filepath.Join(mypkging.current.Dir, ".pkgtest")
	mypkging.paths.Current = mypkging.current

	suite, err := pkgtest.NewSuite(mypkging)
	r.NoError(err)

	suite.Test(t)
}
