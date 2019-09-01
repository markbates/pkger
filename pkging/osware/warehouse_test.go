package hdware

import (
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging/waretest"
	"github.com/stretchr/testify/require"
)

func Test_Warehouse(t *testing.T) {
	r := require.New(t)

	mypkging, err := New()
	r.NoError(err)

	mypkging.current.Dir = filepath.Join(mypkging.current.Dir, ".waretest")
	mypkging.paths.Current = mypkging.current

	suite, err := waretest.NewSuite(mypkging)
	r.NoError(err)

	suite.Test(t)
}
