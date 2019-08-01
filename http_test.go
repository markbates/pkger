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

func Test_HTTP_File_Memory(t *testing.T) {
	r := require.New(t)

	// f, err := createFile(rootIndex, "/cmd/pkger/main.go")
	// r.NoError(err)
	// defer f.Close()
	//
	f, err := Create("/suit/case.txt")
	r.NoError(err)
	f.Write([]byte(radio))
	r.NoError(f.Close())

	dir, err := Open("/")
	r.NoError(err)
	defer dir.Close()

	ts := httptest.NewServer(http.FileServer(dir))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/suit/case.txt")
	r.NoError(err)
	r.Equal(200, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	r.NoError(err)
	r.Contains(string(b), "I wanna bite the hand that feeds me")

}
