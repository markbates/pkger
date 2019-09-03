package pkgtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

const mould = "/easy/listening/sugar.file"
const hart = "/easy/listening/grant.hart"
const husker = "github.com/husker/du"

type Suite struct {
	pkging.Pkger
}

func NewSuite(yourpkging pkging.Pkger) (Suite, error) {
	suite := Suite{
		Pkger: yourpkging,
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
	name = fmt.Sprintf("%T_%s", s.Pkger, name)
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
	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in string
	}{
		{in: mould},
		{in: ":" + mould},
		{in: ip + ":" + mould},
		{in: filepath.Dir(mould)},
		{in: ":" + filepath.Dir(mould)},
		{in: ip + ":" + filepath.Dir(mould)},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))
			r.NoError(s.MkdirAll(filepath.Dir(pt.Name), 0755))

			f, err := s.Create(pt.Name)
			r.NoError(err)
			r.Equal(pt.Name, f.Name())

			fi, err := f.Stat()
			r.NoError(err)

			r.Equal(pt.Name, fi.Name())
			r.Equal(os.FileMode(0644), fi.Mode())
			r.NotZero(fi.ModTime())
			r.NoError(s.RemoveAll(pt.String()))
		})
	}
}

func (s Suite) Test_Create_No_MkdirAll(t *testing.T) {
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
		{in: filepath.Dir(mould)},
		{in: ":" + filepath.Dir(mould)},
		{in: ip + ":" + filepath.Dir(mould)},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))

			_, err = s.Create(pt.Name)
			r.Error(err)
		})
	}
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
		{in: filepath.Dir(mould)},
		{in: ":" + filepath.Dir(mould)},
		{in: ip + ":" + filepath.Dir(mould)},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))

			dir := filepath.Dir(pt.Name)
			r.NoError(s.MkdirAll(dir, 0755))

			fi, err := s.Stat(dir)
			r.NoError(err)

			r.Equal(dir, fi.Name())
			r.Equal(os.FileMode(0755), fi.Mode().Perm())
			r.NotZero(fi.ModTime())
			r.NoError(s.RemoveAll(pt.String()))
		})
	}
}

func (s Suite) Test_Open_File(t *testing.T) {
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
			r.NoError(s.MkdirAll(filepath.Dir(pt.Name), 0755))

			body := "!" + pt.String()

			pkgutil.WriteFile(s, tt.in, []byte(body), 0644)

			f, err := s.Open(tt.in)
			r.NoError(err)

			r.Equal(pt.Name, f.Path().Name)
			b, err := ioutil.ReadAll(f)
			r.NoError(err)
			r.Equal(body, string(b))

			b, err = pkgutil.ReadFile(s, tt.in)
			r.NoError(err)
			r.Equal(body, string(b))

			r.NoError(f.Close())
		})
	}
}

func (s Suite) Test_Parse(t *testing.T) {
	r := require.New(t)

	cur, err := s.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in  string
		exp pkging.Path
	}{
		{in: mould, exp: pkging.Path{Pkg: ip, Name: mould}},
		{in: filepath.Join(cur.Dir, mould), exp: pkging.Path{Pkg: ip, Name: mould}},
		{in: ":" + mould, exp: pkging.Path{Pkg: ip, Name: mould}},
		{in: ip + ":" + mould, exp: pkging.Path{Pkg: ip, Name: mould}},
		{in: ip, exp: pkging.Path{Pkg: ip, Name: "/"}},
		{in: ":", exp: pkging.Path{Pkg: ip, Name: "/"}},
		{in: husker + ":" + mould, exp: pkging.Path{Pkg: husker, Name: mould}},
		{in: husker, exp: pkging.Path{Pkg: husker, Name: "/"}},
		{in: husker + ":", exp: pkging.Path{Pkg: husker, Name: "/"}},
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
			r.NoError(s.MkdirAll(filepath.Dir(pt.Name), 0755))

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

func (s Suite) Test_Remove(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_HTTP_Open(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_HTTP_Readdir(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_ReadFile(t *testing.T) {
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
			r.NoError(s.MkdirAll(filepath.Dir(pt.Name), 0755))

			f, err := s.Create(tt.in)
			r.NoError(err)

			body := "!" + pt.String()
			_, err = io.Copy(f, strings.NewReader(body))
			r.NoError(err)
			r.NoError(f.Close())

			b, err := pkgutil.ReadFile(s, tt.in)
			r.NoError(err)
			r.Equal(body, string(b))
		})
	}
}
