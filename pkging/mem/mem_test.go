package mem

import (
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Pkger(t *testing.T) {
	r := require.New(t)

	info, err := here.Current()
	r.NoError(err)
	r.NotZero(info)

	wh, err := New(info)
	r.NoError(err)

	WithInfo(wh, info)

	suite, err := pkgtest.NewSuite(wh)
	r.NoError(err)

	suite.Test(t)
}
