package pkgtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

const mould = "/easy/listening/sugar.file"
const hart = "/easy/listening/grant.hart"
const husker = "github.com/husker/du"

type Suite struct {
	Name string
	gen  func() (pkging.Pkger, error)
}

func (s Suite) Make() (pkging.Pkger, error) {
	if s.gen == nil {
		return nil, fmt.Errorf("missing generator function")
	}
	return s.gen()
}

func NewSuite(name string, fn func() (pkging.Pkger, error)) (Suite, error) {
	suite := Suite{
		Name: name,
		gen:  fn,
	}
	return suite, nil
}

func (s Suite) Test(t *testing.T) {
	rv := reflect.ValueOf(s)
	rt := rv.Type()
	if rt.NumMethod() == 0 {
		t.Fatalf("something went wrong wrong with %s", s.Name)
	}
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		if !strings.HasPrefix(m.Name, "Test_") {
			continue
		}

		s.sub(t, m)
	}
}

// func (s Suite) clone() (Suite, error) {
// 	if ns, ok := s.Pkger.(Newable); ok {
// 		pkg, err := ns.New()
// 		if err != nil {
// 			return s, err
// 		}
// 		s, err = NewSuite(pkg)
// 		if err != nil {
// 			return s, err
// 		}
// 	}
// 	if ns, ok := s.Pkger.(WithRootable); ok {
// 		dir, err := ioutil.TempDir("")
// 		if err != nil {
// 			return s, err
// 		}
// 		// defer opkg.RemoveAll(dir)
//
// 		pkg, err := ns.WithRoot(dir)
// 		if err != nil {
// 			return s, err
// 		}
// 		s, err = NewSuite(pkg)
// 		if err != nil {
// 			return s, err
// 		}
// 	}
// 	return s, nil
// }

func (s Suite) Run(t *testing.T, name string, fn func(t *testing.T)) {
	t.Run(name, func(st *testing.T) {
		fn(st)
	})
}

func (s Suite) sub(t *testing.T, m reflect.Method) {
	name := fmt.Sprintf("%s/%s", s.Name, m.Name)
	// s, err := s.clone()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	s.Run(t, name, func(st *testing.T) {
		m.Func.Call([]reflect.Value{
			reflect.ValueOf(s),
			reflect.ValueOf(st),
		})
	})
}

func (s Suite) Test_Create(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in string
	}{
		{in: mould},
		{in: ":" + mould},
		{in: ip + ":" + mould},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			r.NoError(pkg.MkdirAll(filepath.Dir(pt.Name), 0755))

			f, err := pkg.Create(pt.Name)
			r.NoError(err)
			r.Equal(pt.Name, f.Name())

			fi, err := f.Stat()
			r.NoError(err)
			r.NoError(f.Close())

			r.Equal(pt.Name, fi.Name())
			r.NotZero(fi.ModTime())
			r.NoError(pkg.RemoveAll(pt.String()))
		})
	}
}

func (s Suite) Test_Create_No_MkdirAll(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
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
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			_, err = pkg.Create(pt.Name)
			r.Error(err)
		})
	}
}

func (s Suite) Test_Current(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	info, err := pkg.Current()
	r.NoError(err)
	r.NotZero(info)
}

func (s Suite) Test_Info(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)

	info, err := pkg.Info(cur.ImportPath)
	r.NoError(err)
	r.NotZero(info)

}

func (s Suite) Test_MkdirAll(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
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
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			dir := filepath.Dir(pt.Name)
			r.NoError(pkg.MkdirAll(dir, 0755))

			fi, err := pkg.Stat(dir)
			r.NoError(err)

			if runtime.GOOS == "windows" {
				dir = strings.Replace(dir, "\\", "/", -1)
			}
			r.Equal(dir, fi.Name())
			r.NotZero(fi.ModTime())
			r.NoError(pkg.RemoveAll(pt.String()))
		})
	}
}

func (s Suite) Test_Open_File(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
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
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			r.NoError(pkg.RemoveAll(pt.String()))
			r.NoError(pkg.MkdirAll(filepath.Dir(pt.Name), 0755))

			body := "!" + pt.String()

			pkgutil.WriteFile(pkg, tt.in, []byte(body), 0644)

			f, err := pkg.Open(tt.in)
			r.NoError(err)

			r.Equal(pt.Name, f.Path().Name)
			b, err := ioutil.ReadAll(f)
			r.NoError(err)
			r.Equal(body, string(b))

			b, err = pkgutil.ReadFile(pkg, tt.in)
			r.NoError(err)
			r.Equal(body, string(b))

			r.NoError(f.Close())
		})
	}
}

func (s Suite) Test_Parse(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in  string
		exp here.Path
	}{
		{in: mould, exp: here.Path{Pkg: ip, Name: mould}},
		{in: filepath.Join(cur.Dir, mould), exp: here.Path{Pkg: ip, Name: mould}},
		{in: ":" + mould, exp: here.Path{Pkg: ip, Name: mould}},
		{in: ip + ":" + mould, exp: here.Path{Pkg: ip, Name: mould}},
		{in: ip, exp: here.Path{Pkg: ip, Name: "/"}},
		{in: ":", exp: here.Path{Pkg: ip, Name: "/"}},
		{in: husker + ":" + mould, exp: here.Path{Pkg: husker, Name: mould}},
		{in: husker, exp: here.Path{Pkg: husker, Name: "/"}},
		{in: husker + ":", exp: here.Path{Pkg: husker, Name: "/"}},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)
			r.Equal(tt.exp, pt)
		})
	}
}

func (s Suite) Test_Stat_Error(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
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
		s.Run(t, tt.in, func(st *testing.T) {

			r := require.New(st)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			r.NoError(pkg.RemoveAll(pt.String()))

			_, err = pkg.Stat(tt.in)
			r.Error(err)
		})
	}
}

func (s Suite) Test_Stat_Dir(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
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
		s.Run(t, tt.in, func(st *testing.T) {

			r := require.New(st)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			r.NoError(pkg.RemoveAll(pt.String()))

			r.NoError(pkg.MkdirAll(pt.Name, 0755))
			info, err := pkg.Stat(tt.in)
			r.NoError(err)
			r.Equal(pt.Name, info.Name())
		})
	}
}

func (s Suite) Test_Stat_File(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
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
		s.Run(t, tt.in, func(st *testing.T) {

			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			r.NoError(pkg.RemoveAll(pt.String()))
			r.NoError(pkg.MkdirAll(filepath.Dir(pt.Name), 0755))

			f, err := pkg.Create(tt.in)
			r.NoError(err)

			_, err = io.Copy(f, strings.NewReader("!"+pt.String()))
			r.NoError(err)
			r.NoError(f.Close())

			info, err := pkg.Stat(tt.in)
			r.NoError(err)
			r.Equal(pt.Name, info.Name())
		})
	}
}

func (s Suite) Test_Walk(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)
	r.NoError(s.LoadFolder(pkg))

	cur, err := pkg.Current()
	r.NoError(err)

	ip := cur.ImportPath

	table := []struct {
		in string
	}{
		{in: ip},
		{in: "/"},
		{in: ":/"},
		{in: ip + ":/"},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)
			var act []string
			err := pkg.Walk(tt.in, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				act = append(act, path)
				return nil
			})
			r.NoError(err)

			exp := []string{
				"github.com/markbates/pkger:/",
				"github.com/markbates/pkger:/main.go",
				"github.com/markbates/pkger:/public",
				"github.com/markbates/pkger:/public/images",
				"github.com/markbates/pkger:/public/images/mark.png",
				"github.com/markbates/pkger:/public/index.html",
				"github.com/markbates/pkger:/templates",
				"github.com/markbates/pkger:/templates/a.txt",
				"github.com/markbates/pkger:/templates/b",
				"github.com/markbates/pkger:/templates/b/b.txt",
			}
			r.Equal(exp, act)
		})
	}

}

func (s Suite) Test_Remove(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)

	ip := cur.ImportPath

	table := []struct {
		in string
	}{
		{in: "/public/images/mark.png"},
		{in: ":/public/images/mark.png"},
		{in: ip + ":/public/images/mark.png"},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)
			r.NoError(s.LoadFolder(pkg))

			_, err = pkg.Stat(tt.in)
			r.NoError(err)

			r.NoError(pkg.Remove(tt.in))

			_, err = pkg.Stat(tt.in)
			r.Error(err)

			r.Error(pkg.Remove("unknown"))
		})
	}

}

func (s Suite) Test_RemoveAll(t *testing.T) {
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
		{in: ":/public"},
		{in: ip + ":/public"},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)
			r.NoError(s.LoadFolder(pkg))

			_, err = pkg.Stat(tt.in)
			r.NoError(err)

			r.NoError(pkg.RemoveAll(tt.in))

			_, err = pkg.Stat(tt.in)
			r.Error(err)
		})
	}

}
