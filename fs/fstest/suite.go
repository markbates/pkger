package fstest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/markbates/pkger/fs"
	"github.com/stretchr/testify/require"
)

type FileSystem struct {
	fs.FileSystem
}

func NewFileSystem(yourfs fs.FileSystem) (*FileSystem, error) {
	suite := &FileSystem{
		FileSystem: yourfs,
	}
	return suite, nil
}

func (s *FileSystem) Test(t *testing.T) {
	rv := reflect.ValueOf(s)
	rt := rv.Type()
	if rt.NumMethod() == 0 {
		t.Fatalf("something went wrong wrong with %s %T", s, s)
	}
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		if !strings.HasPrefix(m.Name, "Test_") {
			continue
		}
		s.sub(t, m)
	}
}

func (s *FileSystem) sub(t *testing.T, m reflect.Method) {
	name := strings.TrimPrefix(m.Name, "Test_")
	name = fmt.Sprintf("%T_%s", s.FileSystem, name)
	t.Run(name, func(st *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				st.Fatal(err)
			}
		}()

		cleaner := func() {
			if err := s.Clean(); err != nil {
				st.Fatal(err)
			}
		}
		cleaner()

		defer cleaner()
		m.Func.Call([]reflect.Value{
			reflect.ValueOf(s),
			reflect.ValueOf(st),
		})
	})
}

func (s *FileSystem) Clean() error {
	pt, err := s.Parse("/")
	if err != nil {
		return err
	}

	_ = pt
	if err := s.RemoveAll(pt.Name); err != nil {
		return err
	}
	//
	// if _, err := s.Stat(pt.Name); err == nil {
	// 	return fmt.Errorf("expected %q to be, you know, not there any more", pt)
	// }
	return nil
}

func (s *FileSystem) Test_Create(t *testing.T) {
	r := require.New(t)

	pt, err := s.Parse("/i/want/candy.song")
	r.NoError(err)

	f, err := s.Create(pt.Name)
	r.NoError(err)
	r.Equal(pt.Name, f.Name())

	fi, err := f.Stat()
	r.NoError(err)

	r.Equal(pt.Name, fi.Name())
	r.Equal(os.FileMode(0644), fi.Mode())
	r.NotZero(fi.ModTime())
	r.NoError(s.RemoveAll(pt.String()))
}

func (s *FileSystem) Test_Current(t *testing.T) {
	r := require.New(t)

	info, err := s.Current()
	r.NoError(err)
	r.NotZero(info)
}

func (s *FileSystem) Test_Info(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	info, err := s.Info(cur.ImportPath)
	r.NoError(err)
	r.NotZero(info)

}

func (s *FileSystem) Test_MkdirAll(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Open(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Parse(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in  string
		exp fs.Path
	}{
		{in: "/foo.go", exp: fs.Path{Pkg: ip, Name: "/foo.go"}},
		{in: filepath.Join(cur.Dir, "foo.go"), exp: fs.Path{Pkg: ip, Name: "/foo.go"}},
		{in: ":/foo.go", exp: fs.Path{Pkg: ip, Name: "/foo.go"}},
		{in: ip + ":/foo.go", exp: fs.Path{Pkg: ip, Name: "/foo.go"}},
		{in: ip, exp: fs.Path{Pkg: ip, Name: "/"}},
		{in: ":", exp: fs.Path{Pkg: ip, Name: "/"}},
		{in: "github.com/old/97s:/foo.go", exp: fs.Path{Pkg: "github.com/old/97s", Name: "/foo.go"}},
		{in: "github.com/old/97s", exp: fs.Path{Pkg: "github.com/old/97s", Name: "/"}},
		{in: "github.com/old/97s:", exp: fs.Path{Pkg: "github.com/old/97s", Name: "/"}},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)
			r.Equal(tt.exp, pt)
		})
	}
}

func (s *FileSystem) Test_ReadFile(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Stat(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in  string
		err bool
	}{
		{in: "/foo.go", err: false},
		{in: ":/foo.go", err: false},
		{in: ip + ":/foo.go", err: false},
		{in: ip, err: false},
		{in: "/no.go", err: true},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			if tt.err {
				_, err := s.Stat(tt.in)
				r.Error(err)
				return
			}

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.Fail(pt.String())
			f, err := s.Create(tt.in)
			r.NoError(err)
			_, err = io.Copy(f, strings.NewReader("!"+pt.String()))
			r.NoError(err)
			r.NoError(f.Close())

			info, err := s.Stat(tt.in)
			r.NoError(err)
			r.Equal(pt.Name, info.Name())
		})
	}
}

func (s *FileSystem) Test_Walk(t *testing.T) {
	panic("not implemented")
}
