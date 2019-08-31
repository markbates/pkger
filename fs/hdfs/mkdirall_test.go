package hdfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MkdirAll(t *testing.T) {
	r := require.New(t)

	fs := NewFS()
	err := fs.MkdirAll("/foo/bar/baz", 0755)
	r.NoError(err)

	fi, err := fs.Stat("/foo/bar/baz")
	r.NoError(err)

	r.Equal("/foo/bar/baz", fi.Name())
	r.Equal(os.FileMode(0755), fi.Mode())
	r.True(fi.IsDir())
}
