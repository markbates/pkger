package hdfs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Stat(t *testing.T) {
	r := require.New(t)

	fs := NewFS()

	_, err := fs.Stat("/i.dont.exist")
	r.Error(err)

	f, err := fs.Create("/i.exist")
	r.NoError(err)
	r.NoError(f.Close())

	fi, err := fs.Stat("/i.exist")
	r.NoError(err)
	r.Equal("/i.exist", fi.Name())
}
