package stdos

import (
	"net/http"
	"os"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

var _ pkging.File = &File{}

type File struct {
	*os.File
	info   *pkging.FileInfo
	her    here.Info
	path   pkging.Path
	pkging pkging.Pkger
}

func NewFile(fx pkging.Pkger, osf *os.File) (*File, error) {

	pt, err := fx.Parse(osf.Name())
	if err != nil {
		return nil, err
	}

	info, err := osf.Stat()
	if err != nil {
		return nil, err
	}

	f := &File{
		File:   osf,
		path:   pt,
		pkging: fx,
	}
	f.info = pkging.WithName(pt.Name, info)

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
	return f.pkging.AbsPath(f.path)
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

func (f *File) Path() pkging.Path {
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
	f.info = pkging.NewFileInfo(info)
	return info, nil
}
