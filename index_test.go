package pkger

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/markbates/pkger/paths"
	"github.com/stretchr/testify/require"
)

func Test_index_Create(t *testing.T) {
	r := require.New(t)

	i := newIndex()

	f, err := i.Create(paths.Path{
		Name: "/hello.txt",
	})
	r.NoError(err)
	r.NotNil(f)

	fi, err := f.Stat()
	r.NoError(err)

	r.Equal("hello.txt", fi.Name())
	r.Equal(os.FileMode(0666), fi.Mode())
	r.NotZero(fi.ModTime())

	her := f.her
	r.NotZero(her)
	r.Equal("github.com/markbates/pkger", her.ImportPath)
}

func Test_index_Create_Write(t *testing.T) {
	r := require.New(t)

	i := newIndex()

	f, err := i.Create(paths.Path{
		Name: "/hello.txt",
	})
	r.NoError(err)
	r.NotNil(f)

	fi, err := f.Stat()
	r.NoError(err)
	r.Zero(fi.Size())

	r.Equal("hello.txt", fi.Name())

	mt := fi.ModTime()
	r.NotZero(mt)

	sz, err := io.Copy(f, strings.NewReader(radio))
	r.NoError(err)
	r.Equal(int64(1381), sz)

	r.NoError(f.Close())
	r.Equal(int64(1381), fi.Size())
	r.NotZero(fi.ModTime())
	r.NotEqual(mt, fi.ModTime())
}
