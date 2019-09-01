package memware

import (
	"testing"

	"github.com/markbates/pkger/fs/fstest"
	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_Warehouse(t *testing.T) {
	r := require.New(t)

	info, err := here.Current()
	r.NoError(err)
	r.NotZero(info)

	myfs, err := New(info)
	r.NoError(err)

	WithInfo(myfs, info)

	suite, err := fstest.NewSuite(myfs)
	r.NoError(err)

	suite.Test(t)
}
