package hdware

import (
	"net/http"
	"os"

	"github.com/markbates/pkger/fs"
	"github.com/markbates/pkger/here"
)

var _ fs.File = &File{}

type File struct {
	*os.File
	info *fs.FileInfo
	her  here.Info
	path fs.Path
	fs   fs.Warehouse
}

func NewFile(fx fs.Warehouse, osf *os.File) (*File, error) {

	pt, err := fx.Parse(osf.Name())
	if err != nil {
		return nil, err
	}

	info, err := osf.Stat()
	if err != nil {
		return nil, err
	}

	f := &File{
		File: osf,
		path: pt,
		fs:   fx,
	}
	f.info = fs.WithName(pt.Name, info)

	her, err := here.Package(pt.Pkg)
	if err != nil {
		return nil, err
	}
	f.her = her
	return f, nil
}

func (f *File) Close() error {
	return f.File.Close()
}

func (f *File) Abs() (string, error) {
	return f.fs.AbsPath(f.path)
}

func (f *File) Info() here.Info {
	return f.her
}

func (f *File) Name() string {
	return f.info.Name()
}

func (f *File) Open(name string) (http.File, error) {
	return f.File, nil
}

func (f *File) Path() fs.Path {
	return f.path
}

func (f *File) Stat() (os.FileInfo, error) {
	if f.info != nil {
		return f.info, nil
	}

	abs, err := f.Abs()
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	f.info = fs.NewFileInfo(info)
	return info, nil
}
