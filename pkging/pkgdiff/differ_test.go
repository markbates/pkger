package pkgdiff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/stretchr/testify/require"
)

func Test_Differ_File(t *testing.T) {
	r := require.New(t)

	const x = "/public/index.html"

	pwd, err := os.Getwd()
	r.NoError(err)

	ch := filepath.Join(pwd,
		"..",
		"..",
		"examples",
		"app")

	os.RemoveAll(filepath.Join(ch, "pkged.go"))

	info := here.Info{
		Dir:        ch,
		ImportPath: "github.com/markbates/pkger/examples/app",
	}
	here.Cache(info.ImportPath, func(s string) (here.Info, error) {
		return info, nil
	})

	disk, err := stdos.New(info)
	r.NoError(err)

	diff := Differ{
		A: disk,
	}

	res, err := diff.File(x)
	r.NoError(err)
	r.Empty(res)
}

func Test_Differ_Dir(t *testing.T) {
	r := require.New(t)

	const x = "/public"

	pwd, err := os.Getwd()
	r.NoError(err)

	ch := filepath.Join(pwd,
		"..",
		"..",
		"examples",
		"app")

	os.RemoveAll(filepath.Join(ch, "pkged.go"))

	info := here.Info{
		Dir:        ch,
		ImportPath: "github.com/markbates/pkger/examples/app",
	}
	here.Cache(info.ImportPath, func(s string) (here.Info, error) {
		return info, nil
	})

	disk, err := stdos.New(info)
	r.NoError(err)

	diff := Differ{
		A: disk,
	}

	res, err := diff.Dir(x)
	r.NoError(err)
	r.Empty(res)
}
