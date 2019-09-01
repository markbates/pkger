package hdware

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/fs"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/maps"
)

var _ fs.Warehouse = &Warehouse{}

type Warehouse struct {
	infos   *maps.Infos
	paths   *maps.Paths
	current here.Info
}

func (f *Warehouse) Abs(p string) (string, error) {
	pt, err := f.Parse(p)
	if err != nil {
		return "", err
	}
	return f.AbsPath(pt)
}

func (f *Warehouse) AbsPath(pt fs.Path) (string, error) {
	if pt.Pkg == f.current.ImportPath {
		return filepath.Join(f.current.Dir, pt.Name), nil
	}
	info, err := f.Info(pt.Pkg)
	if err != nil {
		return "", err
	}
	return filepath.Join(info.Dir, pt.Name), nil
}

func New() (*Warehouse, error) {
	info, err := here.Current()
	if err != nil {
		return nil, err
	}
	return &Warehouse{
		infos: &maps.Infos{},
		paths: &maps.Paths{
			Current: info,
		},
		current: info,
	}, nil
}

func (fx *Warehouse) Create(name string) (fs.File, error) {
	name, err := fx.Abs(name)
	if err != nil {
		return nil, err
	}
	if err := fx.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return nil, err
	}
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return NewFile(fx, f)
}

func (f *Warehouse) Current() (here.Info, error) {
	return f.current, nil
}

func (f *Warehouse) Info(p string) (here.Info, error) {
	info, ok := f.infos.Load(p)
	if ok {
		return info, nil
	}

	info, err := here.Package(p)
	if err != nil {
		return info, err
	}
	f.infos.Store(p, info)
	return info, nil
}

func (f *Warehouse) MkdirAll(p string, perm os.FileMode) error {
	p, err := f.Abs(p)
	if err != nil {
		return err
	}
	return os.MkdirAll(p, perm)
}

func (fx *Warehouse) Open(name string) (fs.File, error) {
	name, err := fx.Abs(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return NewFile(fx, f)
}

func (f *Warehouse) Parse(p string) (fs.Path, error) {
	return f.paths.Parse(p)
}

func (f *Warehouse) ReadFile(s string) ([]byte, error) {
	s, err := f.Abs(s)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(s)
}

func (f *Warehouse) Stat(name string) (os.FileInfo, error) {
	pt, err := f.Parse(name)
	if err != nil {
		return nil, err
	}

	abs, err := f.AbsPath(pt)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}

	info = fs.WithName(pt.Name, fs.NewFileInfo(info))

	return info, nil
}

func (f *Warehouse) Walk(p string, wf filepath.WalkFunc) error {
	fp, err := f.Abs(p)
	if err != nil {
		return err
	}

	pt, err := f.Parse(p)
	if err != nil {
		return err
	}
	err = filepath.Walk(fp, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, fp)
		pt, err := f.Parse(fmt.Sprintf("%s:%s", pt.Pkg, path))
		if err != nil {
			return err
		}
		return wf(pt.String(), fs.WithName(path, fs.NewFileInfo(fi)), nil)
	})

	return err
}

func (fx *Warehouse) Remove(name string) error {
	name, err := fx.Abs(name)
	if err != nil {
		return err
	}
	return os.Remove(name)
}

func (fx *Warehouse) RemoveAll(name string) error {
	name, err := fx.Abs(name)
	if err != nil {
		return err
	}
	return os.RemoveAll(name)
}
