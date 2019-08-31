package hdfs

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Open(t *testing.T) {
	r := require.New(t)

	fs := NewFS()

	_, err := fs.Open("/i.dont.exist")
	r.Error(err)

	f, err := fs.Create("/i.exist")
	r.NoError(err)
	_, err = io.Copy(f, strings.NewReader(radio))
	r.NoError(err)
	r.NoError(f.Close())

	f, err = fs.Open("/i.exist")
	r.NoError(err)
	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.NoError(f.Close())
	r.Equal([]byte(radio), b)
}
