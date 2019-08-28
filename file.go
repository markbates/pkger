package pkger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/markbates/pkger/here"
)

const timeFmt = time.RFC3339Nano

var _ http.File = &File{}

type File struct {
	info     *FileInfo
	her      here.Info
	path     Path
	data     []byte
	parent   Path
	writer   *bytes.Buffer
	reader   io.Reader
	excludes []string
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if sk, ok := f.reader.(io.Seeker); ok {
		return sk.Seek(offset, whence)
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
	fi.size = int64(len(f.data))
	fi.modTime = time.Now()
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

	of, err := f.her.Open(f.FilePath())
	if err != nil {
		return 0, err
	}
	f.reader = of
	return f.reader.Read(p)
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

func (f File) FilePath() string {
	return f.her.FilePath(f.Name())
}

func (f File) Path() Path {
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
	err := Walk(f.Name(), func(pt Path, info os.FileInfo) error {
		if count > 0 && len(infos) == count {
			return io.EOF
		}

		for _, x := range f.excludes {
			rx, err := regexp.Compile(x)
			if err != nil {
				return err
			}
			if rx.MatchString(pt.Name) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if pt.Name == f.parent.Name {
			return nil
		}
		if f.parent.Name != "/" {
			info = WithName(strings.TrimPrefix(info.Name(), f.parent.Name), info)
		}
		infos = append(infos, info)
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
