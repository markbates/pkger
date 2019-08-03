package pkger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_index_Create(t *testing.T) {
	r := require.New(t)

	i := newIndex()

	f, err := i.Create(Path{
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

	f, err := i.Create(Path{
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

func Test_index_JSON(t *testing.T) {
	r := require.New(t)

	i := newIndex()

	f, err := i.Create(Path{
		Name: "/radio.radio",
	})
	r.NoError(err)
	r.NotNil(f)
	fmt.Fprint(f, radio)
	r.NoError(f.Close())

	c, err := i.Stat()
	r.NoError(err)
	r.Equal(curPkg, c.ImportPath)

	_, err = i.Info("github.com/markbates/hepa")
	r.NoError(err)

	r.Equal(1, len(i.Files.Keys()))
	r.Equal(1, len(i.Infos.Keys()))
	r.NotZero(i.Current)

	jason, err := json.Marshal(i)
	r.NoError(err)
	r.NotZero(jason)

	i2 := &index{}

	r.NoError(json.Unmarshal(jason, i2))

	r.NotNil(i2.Infos)
	r.NotNil(i2.Files)
	r.NotZero(i2.Current)
	r.Equal(1, len(i2.Files.Keys()))
	r.Equal(1, len(i2.Infos.Keys()))

	f2, err := i2.Open(Path{Name: "/radio.radio"})
	r.NoError(err)
	r.Equal(f.data, f2.data)
}
