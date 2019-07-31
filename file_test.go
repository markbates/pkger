package pkger

import (
	"io/ioutil"
	"testing"

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
