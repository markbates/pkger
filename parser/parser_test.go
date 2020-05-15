package parser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/stdos"
	"github.com/stretchr/testify/require"
)

func Test_Parser_Ref(t *testing.T) {
	defer func() {
		c := exec.Command("go", "mod", "tidy", "-v")
		c.Run()
	}()
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	disk, err := stdos.New(ref.Info)
	r.NoError(err)

	_, err = pkgtest.LoadFiles("/", ref, disk)
	r.NoError(err)

	res, err := Parse(ref.Info)

	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)
	for _, f := range files {
		fmt.Println(f.Path)
	}
	r.Len(files, 25)
}

func Test_Parser_Ref_Include(t *testing.T) {
	defer func() {
		c := exec.Command("go", "mod", "tidy", "-v")
		c.Run()
	}()
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	disk, err := stdos.New(ref.Info)
	r.NoError(err)

	_, err = pkgtest.LoadFiles("/", ref, disk)
	r.NoError(err)

	res, err := Parse(ref.Info, "github.com/stretchr/testify:/go.mod")

	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)

	l := len(files)
	r.Equal(26, l)
}

func Test_Parser_dotGo_Directory(t *testing.T) {
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	err = os.Mkdir(filepath.Join(ref.Dir, ".go"), 0755)
	r.NoError(err)

	disk, err := stdos.New(ref.Info)
	r.NoError(err)

	_, err = pkgtest.LoadFiles("/", ref, disk)
	r.NoError(err)

	res, err := Parse(ref.Info)
	r.NoError(err)
	r.Equal(11, len(res))
}

func Test_Parser_Example_HTTP(t *testing.T) {
	r := require.New(t)

	cur, err := here.Package("github.com/markbates/pkger")
	r.NoError(err)

	root := filepath.Join(cur.Dir, "examples", "http", "pkger")

	defer func() {
		c := exec.Command("go", "mod", "tidy", "-v")
		c.Run()
	}()

	her, err := here.Dir(root)
	r.NoError(err)

	res, err := Parse(her)
	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)
	r.Len(files, 5)

	for _, f := range files {
		r.True(strings.HasPrefix(f.Abs, her.Dir), "%q %q", f.Abs, her.Dir)
		r.True(strings.HasPrefix(f.Path.Name, "/public"), "%q %q", f.Path.Name, "/public")
	}
}
