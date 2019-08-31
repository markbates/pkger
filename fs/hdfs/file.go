package hdfs

import (
	"net/http"
	"os"

	"github.com/markbates/pkger/fs"
	"github.com/markbates/pkger/here"
)

var _ fs.File = &File{}

type File struct {
	*os.File
	filePath string
	info     *fs.FileInfo
	her      here.Info
	path     fs.Path
	fs       fs.FileSystem
}

func NewFile(fx fs.FileSystem, osf *os.File) (*File, error) {
	info, err := osf.Stat()
	if err != nil {
		return nil, err
	}

	pt, err := fx.Parse(info.Name())
	if err != nil {
		return nil, err
	}

	f := &File{
		File:     osf,
		filePath: info.Name(),
		path:     pt,
		fs:       fx,
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

func (f *File) FilePath() string {
	return f.filePath
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
	if f.info == nil {
		info, err := os.Stat(f.filePath)
		if err != nil {
			return nil, err
		}
		f.info = fs.NewFileInfo(info)
	}
	return f.info, nil
}
