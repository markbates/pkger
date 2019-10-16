package stdos

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_File_Stat_No_Info(t *testing.T) {
	r := require.New(t)

	her, err := here.Current()
	r.NoError(err)
	pkg, err := New(her)
	r.NoError(err)

	f, err := pkg.Open(":/pkging/stdos/file_test.go")
	r.NoError(err)
	defer f.Close()

	sf, ok := f.(*File)
	r.True(ok)

	oi := sf.info
	sf.info = nil

	info, err := sf.Stat()
	r.NoError(err)
	r.Equal(oi.Name(), info.Name())
	// r.Equal("", f.Name())
}

func Test_File_HTTP_Dir(t *testing.T) {
	r := require.New(t)

	her, err := here.Current()
	r.NoError(err)
	pkg, err := New(her)
	r.NoError(err)

	fp := filepath.Join("..", "..", "examples", "app", "public")

	gots := httptest.NewServer(http.FileServer(http.Dir(fp)))
	defer gots.Close()

	dir, err := pkg.Open("/examples/app/public")
	r.NoError(err)

	pkgts := httptest.NewServer(http.FileServer(dir))
	defer pkgts.Close()

	paths := []string{
		"/",
		"/index.html",
		"/images",
		"/images/images/mark.png",
	}

	for _, path := range paths {
		t.Run(path, func(st *testing.T) {
			r := require.New(st)

			gores, err := http.Get(gots.URL + path)
			r.NoError(err)

			pkgres, err := http.Get(pkgts.URL + path)
			r.NoError(err)

			gobody, err := ioutil.ReadAll(gores.Body)
			r.NoError(err)

			pkgbody, err := ioutil.ReadAll(pkgres.Body)
			r.NoError(err)
			r.Equal(string(gobody), string(pkgbody))
		})
	}
}
