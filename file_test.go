package pkger

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/markbates/pkger/paths"
	"github.com/stretchr/testify/require"
)

func Test_File_Open(t *testing.T) {
	r := require.New(t)

	f, err := Open("/file_test.go")
	r.NoError(err)

	r.Equal("file_test.go", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Contains(string(b), "Test_File_Open")
	r.NoError(f.Close())
}

func Test_File_Open_Dir(t *testing.T) {
	r := require.New(t)

	f, err := Open("/cmd")
	r.NoError(err)

	r.Equal("cmd", f.Name())

	r.NoError(f.Close())
}

func Test_File_Read_Memory(t *testing.T) {
	r := require.New(t)

	f, err := Open("/file_test.go")
	r.NoError(err)
	f.data = []byte("hi!")

	r.Equal("file_test.go", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal(string(b), "hi!")
	r.NoError(f.Close())
}

func Test_File_Write(t *testing.T) {
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
