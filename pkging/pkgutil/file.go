package pkgutil

import (
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/markbates/pkger/pkging"
)

type Opener interface {
	Open(name string) (pkging.File, error)
}

type Creator interface {
	Create(name string) (pkging.File, error)
}

type OpenFiler interface {
	OpenFile(name string, flag int, perm os.FileMode) (pkging.File, error)
}

// ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename.
func ReadDir(pkg Opener, dirname string) ([]os.FileInfo, error) {
	f, err := pkg.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

// ReadFile reads the file named by filename and returns the contents. A successful call returns err == nil, not err == EOF. Because ReadFile reads the whole file, it does not treat an EOF from Read as an error to be reported.
func ReadFile(pkg Opener, s string) ([]byte, error) {
	f, err := pkg.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// WriteFile writes data to a file named by filename. If the file does not exist, WriteFile creates it with permissions perm; otherwise WriteFile truncates it before writing.
func WriteFile(pkg Creator, filename string, data []byte, perm os.FileMode) error {
	var f pkging.File
	var err error

	if of, ok := pkg.(OpenFiler); ok {
		f, err = of.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			return err
		}
	}

	if f == nil {
		f, err = pkg.Create(filename)
		if err != nil {
			return err
		}
	}

	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
