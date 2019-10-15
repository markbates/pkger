package pkgtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

// ├── main.go
// ├── public
// │   ├── images
// │   │   ├── mark.png
// │   └── index.html
// └── templates
//     ├── a.txt
//     └── b
//         └── b.txt
var folderFiles = []string{
	"/main.go",
	"/public/images/mark.png",
	"/public/index.html",
	"/templates/a.txt",
	"/templates/b/b.txt",
}

func (s Suite) WriteFolder(path string) error {
	for _, f := range folderFiles {
		f = filepath.Join(path, f)
		if err := os.MkdirAll(filepath.Dir(f), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(f, []byte("!"+f), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (s Suite) LoadFolder(pkg pkging.Pkger) error {
	for _, f := range folderFiles {
		if err := pkg.MkdirAll(filepath.Dir(f), 0755); err != nil {
			return err
		}
		if err := pkgutil.WriteFile(pkg, f, []byte("!"+f), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (s Suite) Test_HTTP_Dir(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)
	ip := cur.ImportPath

	table := []struct {
		in  string
		req string
		exp string
	}{
		{in: "/", req: "/", exp: `>public/</a`},
		{in: ":" + "/", req: "/", exp: `>public/</a`},
		{in: ip + ":" + "/", req: "/", exp: `>public/</a`},
	}

	for _, tt := range table {
		s.Run(t, tt.in+tt.req, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)
			r.NoError(s.LoadFolder(pkg))

			dir, err := pkg.Open(tt.in)
			r.NoError(err)
			defer dir.Close()

			ts := httptest.NewServer(http.FileServer(dir))
			defer ts.Close()

			res, err := http.Get(ts.URL + tt.req)
			r.NoError(err)
			r.Equal(200, res.StatusCode)

			b, err := ioutil.ReadAll(res.Body)
			r.NoError(err)
			r.Contains(string(b), tt.exp)
			r.NotContains(string(b), "mark.png")
		})
	}
}

func (s Suite) Test_HTTP_Dir_IndexHTML(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)
	ip := cur.ImportPath

	table := []struct {
		in  string
		req string
	}{
		{in: "/public", req: "/"},
		{in: ":" + "/public", req: "/"},
		{in: ip + ":" + "/public", req: "/"},
	}

	exp := "index.html"
	for _, tt := range table {
		s.Run(t, tt.in+exp, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			r.NoError(s.LoadFolder(pkg))

			dir, err := pkg.Open(tt.in)
			r.NoError(err)
			defer dir.Close()

			ts := httptest.NewServer(http.FileServer(dir))
			defer ts.Close()

			res, err := http.Get(ts.URL + tt.req)
			r.NoError(err)
			r.Equal(200, res.StatusCode)

			b, err := ioutil.ReadAll(res.Body)
			r.NoError(err)

			body := string(b)
			r.Contains(body, exp)
			r.NotContains(body, "mark.png")
		})
	}
}

func (s Suite) Test_HTTP_File(t *testing.T) {
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
			fmt.Println(">>>TODO pkging/pkgtest/http.go:190: tpub ", tpub)
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
				"/images/mark.png",
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

					exp := strings.ReplaceAll(string(gobody), tdir, "")
					r.Equal(exp, string(pkgbody))
				})
			}
		})
	}

}

//
// func (s Suite) Test_HTTP_Dir_Memory_StripPrefix(t *testing.T) {
// 	r := require.New(t)
//
// 	fs := NewPkger()
// 	r.NoError(Folder.Create(fs))
//
// 	dir, err := fs.Open("/public")
// 	r.NoError(err)
// 	defer dir.Close()
//
// 	ts := httptest.NewServer(http.StripPrefix("/assets/", http.FileServer(dir)))
// 	defer ts.Close()
//
// 	res, err := http.Get(ts.URL + "/assets/images/mark.png")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, _ := ioutil.ReadAll(res.Body)
// 	// r.NoError(err)
// 	r.Contains(string(b), "!/public/images/mark.png")
//
// 	res, err = http.Get(ts.URL + "/assets/images/")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, _ = ioutil.ReadAll(res.Body)
// 	// r.NoError(err)
// 	r.Contains(string(b), `<a href="/mark.png">/mark.png</a>`)
// 	r.NotContains(string(b), `/public`)
// 	r.NotContains(string(b), `/images`)
// 	r.NotContains(string(b), `/go.mod`)
// }
