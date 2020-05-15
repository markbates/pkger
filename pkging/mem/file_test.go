package mem

import (
	"io/ioutil"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_File_Seek(t *testing.T) {
	r := require.New(t)

	info, err := here.Current()
	r.NoError(err)

	pkg, err := New(info)
	r.NoError(err)

	f, err := pkg.Create(":/wilco.band")
	r.NoError(err)

	data := []byte("a shot in the arm")
	f.Write(data)
	r.NoError(f.Close())

	f, err = pkg.Open(":/wilco.band")
	r.NoError(err)

	// seek to end of file before read
	pos, err := f.Seek(0, 2)
	r.NoError(err)
	r.Equal(int64(len(data)), pos)

	// reset seek
	pos, err = f.Seek(0, 0)
	r.NoError(err)
	r.Equal(int64(0), pos)

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal(data, b)

	_, err = f.Seek(0, 0)
	r.NoError(err)

	b, err = ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal(data, b)

	b, err = ioutil.ReadAll(f)
	r.NoError(err)
	r.NotEqual(data, b)

	_, err = f.Seek(10, 0)
	r.NoError(err)

	b, err = ioutil.ReadAll(f)
	r.NoError(err)
	r.NotEqual(data, b)
	r.Equal([]byte("the arm"), b)
}
