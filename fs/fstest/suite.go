package fstest

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/markbates/pkger/fs"
	"github.com/stretchr/testify/require"
)

type FileSystem struct {
	FS fs.FileSystem
}

func NewFileSystem(yourfs fs.FileSystem) (*FileSystem, error) {
	suite := &FileSystem{
		FS: yourfs,
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
	name = fmt.Sprintf("%T_%s", s.FS, name)
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
	name := Path(s.FS)
	err := s.FS.RemoveAll(name)
	if err != nil {
		return err
	}

	if _, err := s.FS.Stat(name); err == nil {
		return fmt.Errorf("expected %q to be, you know, not there any more", name)
	}
	return nil
}

func (s *FileSystem) Test_Create(t *testing.T) {
	r := require.New(t)
	name := Path(s.FS, "i", "want", "candy.song")

	f, err := s.FS.Create(name)
	r.NoError(err)
	r.Equal(name, f.Name())

	fi, err := f.Stat()
	r.NoError(err)

	r.Equal(name, fi.Name())
	r.Equal(os.FileMode(0644), fi.Mode())
	r.NotZero(fi.ModTime())
}

func (s *FileSystem) Test_Current(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Info(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_MkdirAll(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Open(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Parse(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_ReadFile(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Stat(t *testing.T) {
	panic("not implemented")
}

func (s *FileSystem) Test_Walk(t *testing.T) {
	panic("not implemented")
}
