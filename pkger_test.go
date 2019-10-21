package pkger_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	pt, err := pkger.Parse(fmt.Sprintf("%s:/public/index.html", app.Info.ImportPath))

	r.NoError(err)
	r.Equal(app.Info.ImportPath, pt.Pkg)
	r.Equal("/public/index.html", pt.Name)
}

func Test_Abs(t *testing.T) {
	r := require.New(t)

	s, err := pkger.Abs("/rocket.ship")
	r.NoError(err)

	pwd, err := os.Getwd()
	r.NoError(err)
	r.Equal(filepath.Join(pwd, "rocket.ship"), s)
}

func Test_AbsPath(t *testing.T) {
	r := require.New(t)

	s, err := pkger.AbsPath(here.Path{
		Pkg:  "github.com/markbates/pkger",
		Name: "/rocket.ship",
	})
	r.NoError(err)

	pwd, err := os.Getwd()
	r.NoError(err)
	r.Equal(filepath.Join(pwd, "rocket.ship"), s)
}

func Test_Current(t *testing.T) {
	r := require.New(t)

	info, err := pkger.Current()
	r.NoError(err)
	r.Equal("github.com/markbates/pkger", info.ImportPath)
}

func Test_Info(t *testing.T) {
	r := require.New(t)

	info, err := pkger.Info("github.com/markbates/pkger")
	r.NoError(err)
	r.Equal("github.com/markbates/pkger", info.ImportPath)
}

func Test_Create(t *testing.T) {
	r := require.New(t)

	pkger.MkdirAll("/tmp", 0755)
	defer pkger.RemoveAll("/tmp")
	f, err := pkger.Create("/tmp/test.create")
	r.NoError(err)
	r.Equal("/tmp/test.create", f.Name())
	r.NoError(f.Close())
}

func Test_MkdirAll(t *testing.T) {
	r := require.New(t)

	_, err := pkger.Open("/tmp")
	r.Error(err)
	r.NoError(pkger.MkdirAll("/tmp", 0755))
	defer pkger.RemoveAll("/tmp")

	f, err := pkger.Open("/tmp")
	r.NoError(err)
	r.Equal("/tmp", f.Name())
	r.NoError(f.Close())
}

func Test_Stat(t *testing.T) {
	r := require.New(t)

	info, err := pkger.Stat("/go.mod")
	r.NoError(err)
	r.Equal("/go.mod", info.Name())
}

func Test_Walk(t *testing.T) {
	r := require.New(t)

	files := map[string]os.FileInfo{}
	err := pkger.Walk("/pkging/pkgtest/internal/testdata/app", func(path string, info os.FileInfo, err error) error {
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

	pkger.MkdirAll("/tmp", 0755)
	defer pkger.RemoveAll("/tmp")
	f, err := pkger.Create("/tmp/test.create")
	r.NoError(err)
	r.Equal("/tmp/test.create", f.Name())
	r.NoError(f.Close())
	r.NoError(pkger.Remove("/tmp/test.create"))

	_, err = pkger.Stat("/tmp/test.create")
	r.Error(err)
}
