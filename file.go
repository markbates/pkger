package pkger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/here"
)

const timeFmt = time.RFC3339Nano

type File struct {
	info   *FileInfo
	her    here.Info
	path   Path
	data   []byte
	index  *index
	writer io.ReadWriter
	Source io.ReadCloser
}

func (f *File) Close() error {
	defer func() {
		f.Source = nil
		f.writer = nil
	}()
	if f.Source != nil {
		if c, ok := f.Source.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}

	if f.writer == nil {
		return nil
	}

	b, err := ioutil.ReadAll(f.writer)
	if err != nil {
		return err
	}
	f.data = b

	fi := f.info
	fi.size = int64(len(f.data))
	fi.modTime = time.Now()
	f.info = fi
	return nil
}

func (f *File) Read(p []byte) (int, error) {
	if len(f.data) > 0 && len(f.data) <= len(p) {
		return copy(p, f.data), io.EOF
	}

	if len(f.data) > 0 {
		f.Source = ioutil.NopCloser(bytes.NewReader(f.data))
	}

	if f.Source != nil {
		return f.Source.Read(p)
	}

	of, err := f.her.Open(f.Path())
	if err != nil {
		return 0, err
	}
	f.Source = of
	return f.Source.Read(p)
}

func (f *File) Write(b []byte) (int, error) {
	if f.writer == nil {
		f.writer = &bytes.Buffer{}
	}
	i, err := f.writer.Write(b)
	fmt.Println(f.Name(), i, err)
	return i, err
}

func (f File) HereInfo() here.Info {
	return f.her
}

func (f File) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["info"] = f.info
	m["her"] = f.her
	m["path"] = f.path
	m["index"] = f.index
	m["data"] = f.data
	if len(f.data) == 0 {
		b, err := ioutil.ReadAll(&f)
		if err != nil {
			return nil, err
		}
		m["data"] = b
	}
	return json.Marshal(m)
}

func (f *File) UnmarshalJSON(b []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	info, ok := m["info"]
	if !ok {
		return fmt.Errorf("missing info")
	}
	f.info = &FileInfo{}
	if err := json.Unmarshal(info, f.info); err != nil {
		return err
	}

	her, ok := m["her"]
	if !ok {
		return fmt.Errorf("missing her")
	}
	if err := json.Unmarshal(her, &f.her); err != nil {
		return err
	}

	path, ok := m["path"]
	if !ok {
		return fmt.Errorf("missing path")
	}
	if err := json.Unmarshal(path, &f.path); err != nil {
		return err
	}

	ind, ok := m["index"]
	if !ok {
		return fmt.Errorf("missing index")
	}
	f.index = newIndex()
	if err := json.Unmarshal(ind, f.index); err != nil {
		return err
	}
	return nil
}

func (f *File) Open(name string) (http.File, error) {
	if f.index == nil {
		f.index = newIndex()
	}
	pt, err := Parse(name)
	if err != nil {
		return nil, err
	}

	if len(pt.Pkg) == 0 {
		pt.Pkg = f.path.Pkg
	}

	h := httpFile{
		crs: &byteCRS{bytes.NewReader(f.data)},
	}

	if pt == f.path {
		h.File = f
	} else {
		of, err := f.index.Open(pt)
		if err != nil {
			return nil, err
		}
		defer of.Close()
		h.File = of
	}

	if len(f.data) > 0 {
		return h, nil
	}

	bf, err := f.her.Open(h.File.Path())
	if err != nil {
		return h, err
	}

	fi, err := bf.Stat()
	if err != nil {
		return h, err
	}

	if fi.IsDir() {
		return h, nil
	}

	if err != nil {
		return nil, err
	}

	h.crs = bf
	return h, nil
}

func (f File) Stat() (os.FileInfo, error) {
	if f.info == nil {
		return nil, os.ErrNotExist
	}
	return f.info, nil
}

func (f File) Name() string {
	return f.info.Name()
}

func (f File) Path() string {
	return f.her.FilePath(f.Name())
}

func (f File) String() string {
	if f.info == nil {
		return ""
	}
	b, _ := json.MarshalIndent(f.info, "", "  ")
	return string(b)
}

// Readdir reads the contents of the directory associated with file and returns a slice of up to n FileInfo values, as would be returned by Lstat, in directory order. Subsequent calls on the same file will yield further FileInfos.
//
// If n > 0, Readdir returns at most n FileInfo structures. In this case, if Readdir returns an empty slice, it will return a non-nil error explaining why. At the end of a directory, the error is io.EOF.
//
// If n <= 0, Readdir returns all the FileInfo from the directory in a single slice. In this case, if Readdir succeeds (reads all the way to the end of the directory), it returns the slice and a nil error. If it encounters an error before the end of the directory, Readdir returns the FileInfo read until that point and a non-nil error.
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	of, err := f.her.Open(f.Path())
	if err != nil {
		return nil, err
	}
	defer of.Close()
	return of.Readdir(count)
}
