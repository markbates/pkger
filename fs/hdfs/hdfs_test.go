package hdfs

import (
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/fs/fstest"
	"github.com/stretchr/testify/require"
)

func Test_FS(t *testing.T) {
	r := require.New(t)

	myfs, err := New()
	r.NoError(err)

	myfs.current.Dir = filepath.Join(myfs.current.Dir, ".fstest")
	myfs.paths.Current = myfs.current

	suite, err := fstest.NewSuite(myfs)
	r.NoError(err)

	suite.Test(t)
}
