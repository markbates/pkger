package pkger

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_index_Create(t *testing.T) {
	r := require.New(t)

	f, err := Create("/hello.txt")
	r.NoError(err)
	r.NotNil(f)

	fi, err := f.Stat()
	r.NoError(err)

	r.Equal("/hello.txt", fi.Name())
	r.Equal(os.FileMode(0666), fi.Mode())
	r.NotZero(fi.ModTime())

	her := f.her
	r.NotZero(her)
	r.Equal("github.com/markbates/pkger", her.ImportPath)
}

func Test_index_Create_Write(t *testing.T) {
	r := require.New(t)

	f, err := Create("/hello.txt")
	r.NoError(err)
	r.NotNil(f)

	fi, err := f.Stat()
	r.NoError(err)
	r.Zero(fi.Size())

	r.Equal("/hello.txt", fi.Name())

	mt := fi.ModTime()
	r.NotZero(mt)

	sz, err := io.Copy(f, strings.NewReader(radio))
	r.NoError(err)
	r.Equal(int64(1381), sz)

	r.NoError(f.Close())
	r.Equal(int64(1381), fi.Size())
	r.NotZero(fi.ModTime())
	r.NotEqual(mt, fi.ModTime())
}

// TODO
// func Test_index_JSON(t *testing.T) {
// 	r := require.New(t)
//
// 	f, err := Create("/radio.radio")
// 	r.NoError(err)
// 	r.NotNil(f)
// 	fmt.Fprint(f, radio)
// 	r.NoError(f.Close())
//
// 	c, err := Stat()
// 	r.NoError(err)
// 	r.Equal(curPkg, c.ImportPath)
//
// 	_, err = Info("github.com/markbates/hepa")
// 	r.NoError(err)
//
// 	r.Equal(1, len(filesCache.Keys()))
// 	r.Equal(1, len(infosCache.Keys()))
// 	r.NotZero(cur)
//
// 	jason, err := json.Marshal(i)
// 	r.NoError(err)
// 	r.NotZero(jason)
//
// 	i2 := &index{}
//
// 	r.NoError(json.Unmarshal(jason, i2))
//
// 	r.NotNil(i2.infosCache)
// 	r.NotNil(i2.filesCache)
// 	r.NotZero(i2.cur)
// 	r.Equal(1, len(i2.filesCache.Keys()))
// 	r.Equal(1, len(i2.infosCache.Keys()))
//
// 	f2, err := i2.Open(Path{Name: "/radio.radio"})
// 	r.NoError(err)
// 	r.Equal(f.data, f2.data)
// }

func Test_index_Parse(t *testing.T) {
	table := []struct {
		in  string
		out string
	}{
		{in: "", out: curPkg + ":/"},
		{in: curPkg, out: curPkg + ":/"},
		// {in: curPkg + "/foo.go", out: curPkg + ":/foo.go"},
		// {in: "/foo.go", out: curPkg + ":/foo.go"},
		{in: "github.com/markbates/pkger/internal/examples/app", out: "github.com/markbates/pkger/internal/examples/app:/"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)
			pt, err := Parse(tt.in)
			r.NoError(err)
			r.Equal(tt.out, pt.String())
		})
	}
}
