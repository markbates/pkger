package pkgtest

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func (s Suite) WriteFolder(root string) error {
	app, err := App()
	if err != nil {
		return err
	}

	return filepath.Walk(app.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		of, err := os.Open(path)
		if err != nil {
			return err
		}
		defer of.Close()

		path = strings.TrimPrefix(path, app.Dir)
		path = filepath.Join(root, path)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, of); err != nil {
			return err
		}
		return nil
	})
}

func (s Suite) LoadFolder(pkg pkging.Pkger) error {
	app, err := App()
	if err != nil {
		return err
	}

	return filepath.Walk(app.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		of, err := os.Open(path)
		if err != nil {
			return err
		}
		defer of.Close()

		if a, ok := pkg.(pkging.Adder); ok {
			return a.Add(of)
		}

		path = strings.TrimPrefix(path, app.Dir)

		pt, err := pkg.Parse(path)
		if err != nil {
			return err
		}

		if err := pkg.MkdirAll(filepath.Dir(pt.Name), 0755); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		f, err := pkg.Create(pt.String())
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, of); err != nil {
			return err
		}
		return nil
	})
}

func (s Suite) Test_HTTP(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)
	ip := cur.ImportPath

	table := []struct {
		in string
	}{
		{in: "/public"},
		{in: ":" + "/public"},
		{in: ip + ":" + "/public"},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			r.NoError(s.LoadFolder(pkg))

			tdir, err := ioutil.TempDir("", "")
			r.NoError(err)
			defer os.RemoveAll(tdir)
			r.NoError(s.WriteFolder(tdir))

			tpub := filepath.Join(tdir, "public")
			gots := httptest.NewServer(http.FileServer(http.Dir(tpub)))
			defer gots.Close()

			dir, err := pkg.Open(tt.in)
			r.NoError(err)
			defer dir.Close()

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

					exp := strings.ReplaceAll(string(gobody), tdir, "")
					exp = clean(exp)
					r.Equal(exp, clean(string(pkgbody)))
				})
			}
		})
	}
}
