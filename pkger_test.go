package pkger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	r := require.New(t)

	pt, err := Parse("github.com/rocket/ship:/little")
	r.NoError(err)
	r.Equal("github.com/rocket/ship", pt.Pkg)
	r.Equal("/little", pt.Name)
}

func Test_Current(t *testing.T) {
	r := require.New(t)

	info, err := Current()
	r.NoError(err)
	r.Equal("github.com/markbates/pkger", info.ImportPath)
}

func Test_Info(t *testing.T) {
	r := require.New(t)

	info, err := Info("github.com/markbates/pkger")
	r.NoError(err)
	r.Equal("github.com/markbates/pkger", info.ImportPath)
}

func Test_Create(t *testing.T) {
	r := require.New(t)

	MkdirAll("/tmp", 0755)
	defer RemoveAll("/tmp")
	f, err := Create("/tmp/test.create")
	r.NoError(err)
	r.Equal("github.com/markbates/pkger:/tmp/test.create", f.Name())
	r.NoError(f.Close())
}

func Test_MkdirAll(t *testing.T) {
	r := require.New(t)

	_, err := Open("/tmp")
	r.Error(err)
	r.NoError(MkdirAll("/tmp", 0755))
	defer RemoveAll("/tmp")

	f, err := Open("/tmp")
	r.NoError(err)
	r.Equal("github.com/markbates/pkger:/tmp", f.Name())
	r.NoError(f.Close())
}

func Test_Stat(t *testing.T) {
	r := require.New(t)

	info, err := Stat("/go.mod")
	r.NoError(err)
	r.Equal("go.mod", info.Name())
}

func Test_Walk(t *testing.T) {
	r := require.New(t)

	files := map[string]os.FileInfo{}
	err := Walk("/pkging/pkgtest/testdata/ref", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files[path] = info
		return nil
	})
	r.NoError(err)

	r.True(len(files) > 10)
}

func Test_Remove(t *testing.T) {
	r := require.New(t)

	MkdirAll("/tmp", 0755)
	defer RemoveAll("/tmp")
	f, err := Create("/tmp/test.create")
	r.NoError(err)
	r.NoError(f.Close())
	r.NoError(Remove("/tmp/test.create"))

	_, err = Stat("/tmp/test.create")
	r.Error(err)
}
