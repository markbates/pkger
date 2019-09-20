package stdos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/maps"
	"github.com/markbates/pkger/pkging"
)

var _ pkging.Pkger = &Pkger{}

type Pkger struct {
	infos   *maps.Infos
	paths   *maps.Paths
	current here.Info
}

func (f *Pkger) Abs(p string) (string, error) {
	pt, err := f.Parse(p)
	if err != nil {
		return "", err
	}
	return f.AbsPath(pt)
}

func (f *Pkger) AbsPath(pt pkging.Path) (string, error) {
	if pt.Pkg == f.current.ImportPath {
		return filepath.Join(f.current.Dir, pt.Name), nil
	}
	info, err := f.Info(pt.Pkg)
	if err != nil {
		return "", err
	}
	return filepath.Join(info.Dir, pt.Name), nil
}

func New() (*Pkger, error) {
	info, err := here.Current()
	if err != nil {
		return nil, err
	}
	p := &Pkger{
		infos: &maps.Infos{},
		paths: &maps.Paths{
			Current: info,
		},
		current: info,
	}
	p.infos.Store(info.ImportPath, info)
	return p, nil
}

func (fx *Pkger) Create(name string) (pkging.File, error) {
	name, err := fx.Abs(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return NewFile(fx, f)
}

func (f *Pkger) Current() (here.Info, error) {
	return f.current, nil
}

func (f *Pkger) Info(p string) (here.Info, error) {
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

func (f *Pkger) MkdirAll(p string, perm os.FileMode) error {
	p, err := f.Abs(p)
	if err != nil {
		return err
	}
	return os.MkdirAll(p, perm)
}

func (fx *Pkger) Open(name string) (pkging.File, error) {
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

func (f *Pkger) Parse(p string) (pkging.Path, error) {
	return f.paths.Parse(p)
}

func (f *Pkger) Stat(name string) (os.FileInfo, error) {
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

	info = pkging.WithName(pt.Name, pkging.NewFileInfo(info))

	return info, nil
}

func (f *Pkger) Walk(p string, wf filepath.WalkFunc) error {
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
		return wf(pt.String(), pkging.WithName(path, pkging.NewFileInfo(fi)), nil)
	})

	return err
}

func (fx *Pkger) Remove(name string) error {
	name, err := fx.Abs(name)
	if err != nil {
		return err
	}
	return os.Remove(name)
}

func (fx *Pkger) RemoveAll(name string) error {
	name, err := fx.Abs(name)
	if err != nil {
		return err
	}
	return os.RemoveAll(name)
}
