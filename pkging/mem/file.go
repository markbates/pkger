package mem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

const timeFmt = time.RFC3339Nano

var _ pkging.File = &File{}

type File struct {
	info   *pkging.FileInfo
	her    here.Info
	path   pkging.Path
	data   []byte
	parent pkging.Path
	writer *bytes.Buffer
	reader io.Reader
	pkging pkging.Pkger
}

func (f *File) Seek(ofpkginget int64, whence int) (int64, error) {
	if sk, ok := f.reader.(io.Seeker); ok {
		return sk.Seek(ofpkginget, whence)
	}
	return 0, nil
}

func (f *File) Close() error {
	defer func() {
		f.reader = nil
		f.writer = nil
	}()
	if f.reader != nil {
		if c, ok := f.reader.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}

	if f.writer == nil {
		return nil
	}

	f.data = f.writer.Bytes()

	fi := f.info
	fi.Details.Size = int64(len(f.data))
	fi.Details.ModTime = pkging.ModTime(time.Now())
	f.info = fi
	return nil
}

func (f *File) Read(p []byte) (int, error) {
	if len(f.data) > 0 && f.reader == nil {
		f.reader = bytes.NewReader(f.data)
	}

	if f.reader != nil {
		return f.reader.Read(p)
	}

	return 0, fmt.Errorf("unable to read %s", f.Name())
}

func (f *File) Write(b []byte) (int, error) {
	if f.writer == nil {
		f.writer = &bytes.Buffer{}
	}
	i, err := f.writer.Write(b)
	return i, err
}

func (f File) Info() here.Info {
	return f.her
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

func (f File) Abs() (string, error) {
	return f.pkging.AbsPath(f.Path())
}

func (f File) Path() pkging.Path {
	return f.path
}

func (f File) String() string {
	return f.Path().String()
}

func (f File) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		if st.Flag('+') {
			b, err := json.MarshalIndent(f, "", "  ")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				return
			}
			fmt.Fprint(st, string(b))
			return
		}
		fmt.Fprint(st, f.String())
	case 'q':
		fmt.Fprintf(st, "%q", f.String())
	default:
		fmt.Fprint(st, f.String())
	}
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	var infos []os.FileInfo
	root := f.Path().String()
	err := f.pkging.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if count > 0 && len(infos) == count {
			return io.EOF
		}

		pt, err := f.pkging.Parse(path)
		if err != nil {
			return err
		}
		if pt.Name == f.parent.Name {
			return nil
		}

		info = pkging.WithName(strings.TrimPrefix(info.Name(), f.parent.Name), info)
		infos = append(infos, info)
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return infos, nil
		}
		if err != io.EOF {
			return nil, err
		}
	}
	return infos, nil

}

func (f *File) Open(name string) (http.File, error) {
	pt, err := f.pkging.Parse(name)
	if err != nil {
		return nil, err
	}

	if pt == f.path {
		return f, nil
	}

	pt.Name = path.Join(f.Path().Name, pt.Name)

	di, err := f.pkging.Open(pt.String())
	if err != nil {
		return nil, err
	}

	fi, err := di.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		d2 := &File{
			info:   pkging.NewFileInfo(fi),
			her:    di.Info(),
			path:   pt,
			parent: f.path,
			pkging: f.pkging,
		}
		di = d2
	}
	return di, nil
}
