package pkger

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HTTP_File(t *testing.T) {
	r := require.New(t)

	f, err := Open(".")
	r.NoError(err)

	ts := httptest.NewServer(http.FileServer(f))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/cmd/pkger/main.go")
	r.NoError(err)
	r.Equal(200, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	r.NoError(err)
	r.Contains(string(b), "does not compute")

	r.NoError(f.Close())
}

func Test_HTTP_Dir(t *testing.T) {
	r := require.New(t)

	f, err := Open(".")
	r.NoError(err)

	ts := httptest.NewServer(http.FileServer(f))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/cmd/pkger")
	r.NoError(err)
	r.Equal(200, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	r.NoError(err)
	r.Contains(string(b), `<a href="main.go">main.go</a>`)

	r.NoError(f.Close())
}
