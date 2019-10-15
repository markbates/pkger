package stdos

import (
	"net/http"
	"os"
	"path"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

var _ pkging.File = &File{}

type File struct {
	*os.File
	info   *pkging.FileInfo
	her    here.Info
	path   here.Path
	pkging pkging.Pkger
}

// func NewFile(her here.Info, fx pkging.Pkger, osf *os.File) (*File, error) {
// 	// fmt.Println(">>>TODO pkging/stdos/file.go:23: her.ImportPath ", her.ImportPath)
// 	name := osf.Name()
// 	pt, err := her.Parse(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	info, err := osf.Stat()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	f := &File{
// 		File:   osf,
// 		path:   pt,
// 		pkging: fx,
// 		her:    her,
// 	}
// 	f.info = pkging.WithName(pt.Name, info)
//
// 	return f, nil
// }

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

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	osinfos, err := f.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	infos := make([]os.FileInfo, len(osinfos))
	for i, info := range osinfos {
		infos[i] = pkging.WithRelName(info.Name(), info)
	}
	return infos, err
}
func (f *File) Open(name string) (http.File, error) {
	fp := path.Join(f.Path().Name, name)
	f2, err := f.pkging.Open(fp)
	if err != nil {
		return nil, err
	}
	return f2, nil
}

func (f *File) Path() here.Path {
	return f.path
}

func (f *File) Stat() (os.FileInfo, error) {
	if f.info != nil {
		return f.info, nil
	}

	info, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	f.info = pkging.WithName(f.path.Name, info)
	return f.info, nil
}
