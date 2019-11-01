package pkgtest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func HTTPTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	_, err := LoadFiles("/public", ref, pkg)
	r.NoError(err)

	defer os.RemoveAll(ref.Dir)

	fp := filepath.Join(ref.Dir, "public")
	gots := httptest.NewServer(http.FileServer(http.Dir(fp)))
	defer gots.Close()

	dir, err := pkg.Open("/public")
	r.NoError(err)

	pkgts := httptest.NewServer(http.FileServer(dir))
	defer pkgts.Close()

	paths := []string{
		"/",
		"/index.html",
		"/images",
		"/images/img1.png",
	}

	for _, path := range paths {
		t.Run(path, func(st *testing.T) {
			r := require.New(st)

			gores, err := http.Get(gots.URL + path)
			r.NoError(err)

			pkgres, err := http.Get(pkgts.URL + path)
			r.NoError(err)

			r.Equal(gores.StatusCode, pkgres.StatusCode)

			gobody, err := ioutil.ReadAll(gores.Body)
			r.NoError(err)

			pkgbody, err := ioutil.ReadAll(pkgres.Body)
			r.NoError(err)

			exp := string(gobody)
			act := string(pkgbody)
			r.Equal(exp, act)
			// exp := strings.ReplaceAll(string(gobody), tdir, "")
			// exp = clean(exp)
			// r.Equal(exp, clean(string(pkgbody)))
		})
	}
}
