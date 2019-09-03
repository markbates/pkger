package mem

import (
	"fmt"
	"os"
	"strings"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/maps"
	"github.com/markbates/pkger/pkging"
)

var _ pkging.Pkger = &Pkger{}

func WithInfo(fx *Pkger, infos ...here.Info) {
	for _, info := range infos {
		fx.infos.Store(info.ImportPath, info)
	}
}

func New(info here.Info) (*Pkger, error) {
	f := &Pkger{
		infos: &maps.Infos{},
		paths: &maps.Paths{
			Current: info,
		},
		files:   &maps.Files{},
		current: info,
	}
	return f, nil
}

type Pkger struct {
	infos   *maps.Infos
	paths   *maps.Paths
	files   *maps.Files
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
	return pt.String(), nil
}

func (f *Pkger) Current() (here.Info, error) {
	return f.current, nil
}

func (f *Pkger) Info(p string) (here.Info, error) {
	info, ok := f.infos.Load(p)
	if !ok {
		return info, fmt.Errorf("no such package %q", p)
	}

	return info, nil
}

func (f *Pkger) Parse(p string) (pkging.Path, error) {
	return f.paths.Parse(p)
}

func (fx *Pkger) Remove(name string) error {
	pt, err := fx.Parse(name)
	if err != nil {
		return err
	}

	if _, ok := fx.files.Load(pt); !ok {
		return &os.PathError{"remove", pt.String(), fmt.Errorf("no such file or directory")}
	}

	fx.files.Delete(pt)
	return nil
}

func (fx *Pkger) RemoveAll(name string) error {
	pt, err := fx.Parse(name)
	if err != nil {
		return err
	}

	fx.files.Range(func(key pkging.Path, file pkging.File) bool {
		if strings.HasPrefix(key.Name, pt.Name) {
			fx.files.Delete(key)
		}
		return true
	})

	return nil
}
