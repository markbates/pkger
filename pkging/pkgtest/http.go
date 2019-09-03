package pkgtest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

// examples/app
// ├── Dockerfile
// ├── Makefile
// ├── go.mod
// ├── go.sum
// ├── main.go
// ├── public
// │   ├── images
// │   │   ├── mark-small.png
// │   │   ├── mark.png
// │   │   ├── mark_250px.png
// │   │   └── mark_400px.png
// │   └── index.html
// └── templates
//     ├── a.txt
//     └── b
//         └── b.txt

func (s Suite) LoadFolder() error {
	files := []string{
		"/main.go",
		"/public/images/mark.png",
		"/public/index.html",
		"/templates/a.txt",
		"/templates/b/b.txt",
	}

	for _, f := range files {
		if err := s.MkdirAll(filepath.Dir(f), 0755); err != nil {
			return err
		}
		if err := pkgutil.WriteFile(s, f, []byte("!"+f), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (s Suite) Test_HTTP_Dir(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
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

			r := require.New(t)
			r.NoError(s.LoadFolder())

			dir, err := s.Open(tt.in)
			r.NoError(err)
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

	cur, err := s.Current()
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

			r.NoError(s.LoadFolder())

			dir, err := s.Open(tt.in)
			r.NoError(err)

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

// func (s Suite) Test_HTTP_File_Memory(t *testing.T) {
// 	r := require.New(t)
//
// 	fs := NewPkger()
// 	r.NoError(Folder.Create(fs))
//
// 	dir, err := fs.Open("/")
// 	r.NoError(err)
// 	ts := httptest.NewServer(http.FileServer(dir))
// 	defer ts.Close()
//
// 	res, err := http.Get(ts.URL + "/public/images/mark.png")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, err := ioutil.ReadAll(res.Body)
// 	r.NoError(err)
// 	r.Contains(string(b), `!/public/images/mark.png`)
// }
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
