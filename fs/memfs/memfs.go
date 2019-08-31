package memfs

import (
	"io/ioutil"

	"github.com/markbates/pkger/fs"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/maps"
)

var _ fs.FileSystem = &FS{}

func New(info here.Info) (*FS, error) {
	f := &FS{
		infos: &maps.Infos{},
		paths: &maps.Paths{},
		files: &maps.Files{},
	}
	return f, nil
}

type FS struct {
	infos   *maps.Infos
	paths   *maps.Paths
	files   *maps.Files
	current here.Info
}

func (f *FS) Current() (here.Info, error) {
	return f.current, nil
}

func (f *FS) Info(p string) (here.Info, error) {
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

func (f *FS) Parse(p string) (fs.Path, error) {
	return f.paths.Parse(p)
}

func (fx *FS) ReadFile(s string) ([]byte, error) {
	f, err := fx.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
