package hdfs

import (
	"testing"

	"github.com/markbates/pkger/fs/fstest"
	"github.com/stretchr/testify/require"
)

func Test_FS(t *testing.T) {
	r := require.New(t)

	myfs, err := New()
	r.NoError(err)

	suite, err := fstest.NewFileSystem(myfs)
	r.NoError(err)

	suite.Test(t)
}
