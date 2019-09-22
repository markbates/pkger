package pkger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	r := require.New(t)

	pt, err := Parse("github.com/rocket/ship:/little")
	r.NoError(err)
	r.Equal("github.com/rocket/ship", pt.Pkg)
	r.Equal("/little", pt.Name)
}

func Test_Abs(t *testing.T) {
	r := require.New(t)

	s, err := Abs(":/rocket.ship")
	r.NoError(err)

	pwd, err := os.Getwd()
	r.NoError(err)
	r.Equal(filepath.Join(pwd, "rocket.ship"), s)
}

func Test_AbsPath(t *testing.T) {
	r := require.New(t)

	s, err := AbsPath(pkging.Path{
		Pkg:  "github.com/markbates/pkger",
		Name: "/rocket.ship",
	})
	r.NoError(err)

	pwd, err := os.Getwd()
	r.NoError(err)
	r.Equal(filepath.Join(pwd, "rocket.ship"), s)
}
