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

const mould = "/easy/listening/sugar.file"
const hart = "/easy/listening/grant.hart"
const husker = "github.com/husker/du"

type Suite struct {
	fs.Warehouse
}

func NewSuite(yourfs fs.Warehouse) (Suite, error) {
	suite := Suite{
		Warehouse: yourfs,
	}
	return suite, nil
}

func (s Suite) Test(t *testing.T) {
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

func (s Suite) sub(t *testing.T, m reflect.Method) {
	name := strings.TrimPrefix(m.Name, "Test_")
	name = fmt.Sprintf("%T_%s", s.Warehouse, name)
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

func (s Suite) Clean() error {
	pt, err := s.Parse("/")
	if err != nil {
		return err
	}

	_ = pt
	if err := s.RemoveAll(pt.Name); err != nil {
		return err
	}

	if _, err := s.Stat(pt.Name); err == nil {
		return fmt.Errorf("expected %q to be, you know, not there any more", pt)
	}
	return nil
}

func (s Suite) Test_Create(t *testing.T) {
	r := require.New(t)

	pt, err := s.Parse(mould)
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

func (s Suite) Test_Current(t *testing.T) {
	r := require.New(t)

	info, err := s.Current()
	r.NoError(err)
	r.NotZero(info)
}

func (s Suite) Test_Info(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	info, err := s.Info(cur.ImportPath)
	r.NoError(err)
	r.NotZero(info)

}

func (s Suite) Test_MkdirAll(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_Open(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_Parse(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in  string
		exp fs.Path
	}{
		{in: mould, exp: fs.Path{Pkg: ip, Name: mould}},
		{in: filepath.Join(cur.Dir, mould), exp: fs.Path{Pkg: ip, Name: mould}},
		{in: ":" + mould, exp: fs.Path{Pkg: ip, Name: mould}},
		{in: ip + ":" + mould, exp: fs.Path{Pkg: ip, Name: mould}},
		{in: ip, exp: fs.Path{Pkg: ip, Name: "/"}},
		{in: ":", exp: fs.Path{Pkg: ip, Name: "/"}},
		{in: husker + ":" + mould, exp: fs.Path{Pkg: husker, Name: mould}},
		{in: husker, exp: fs.Path{Pkg: husker, Name: "/"}},
		{in: husker + ":", exp: fs.Path{Pkg: husker, Name: "/"}},
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

func (s Suite) Test_ReadFile(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_Stat_Error(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath

	table := []struct {
		in string
	}{
		{in: hart},
		{in: ":" + hart},
		{in: ip},
		{in: ip + ":"},
		{in: ip + ":" + hart},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {

			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))

			_, err = s.Stat(tt.in)
			r.Error(err)
		})
	}
}

func (s Suite) Test_Stat_Dir(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	dir := filepath.Dir(mould)
	ip := cur.ImportPath

	table := []struct {
		in string
	}{
		{in: ip},
		{in: dir},
		{in: ":" + dir},
		{in: ip + ":" + dir},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {

			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))

			r.NoError(s.MkdirAll(pt.Name, 0755))
			info, err := s.Stat(tt.in)
			r.NoError(err)
			r.Equal(pt.Name, info.Name())
		})
	}
}

func (s Suite) Test_Stat_File(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in string
	}{
		{in: mould},
		{in: ":" + mould},
		{in: ip + ":" + mould},
		{in: hart},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {

			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))

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

func (s Suite) Test_Walk(t *testing.T) {
	panic("not implemented")
}
