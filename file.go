package pkger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/debug"
)

const timeFmt = time.RFC3339Nano

var _ http.File = &File{}

type File struct {
	info   *FileInfo
	her    here.Info
	path   Path
	data   []byte
	parent Path
	writer *bytes.Buffer
	reader io.Reader
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

func (f File) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["info"] = f.info
	m["her"] = f.her
	m["path"] = f.path
	m["data"] = f.data
	m["parent"] = f.parent
	if !f.info.virtual {
		if len(f.data) == 0 && !f.info.IsDir() {
			b, err := ioutil.ReadAll(&f)
			if err != nil {
				return nil, err
			}
			m["data"] = b
		}
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

	parent, ok := m["parent"]
	if !ok {
		return fmt.Errorf("missing parent")
	}
	if err := json.Unmarshal(parent, &f.parent); err != nil {
		return err
	}

	if err := json.Unmarshal(m["data"], &f.data); err != nil {
		return err
	}

	return nil
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

// Readdir reads the contents of the directory associated with file and returns a slice of up to n FileInfo values, as would be returned by Lstat, in directory order. Subsequent calls on the same file will yield further FileInfos.
//
// If n > 0, Readdir returns at most n FileInfo structures. In this case, if Readdir returns an empty slice, it will return a non-nil error explaining why. At the end of a directory, the error is io.EOF.
//
// If n <= 0, Readdir returns all the FileInfo from the directory in a single slice. In this case, if Readdir succeeds (reads all the way to the end of the directory), it returns the slice and a nil error. If it encounters an error before the end of the directory, Readdir returns the FileInfo read until that point and a non-nil error.
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	var infos []os.FileInfo
	defer func() {
		fmt.Println(f.Name(), len(infos))
	}()
	err := Walk(f.Name(), func(pt Path, info os.FileInfo) error {
		if count > 0 && len(infos) == count {
			return io.EOF
		}
		debug.Debug("[PKGER] [*file|Readdir] %d %q %q", count, f, pt)

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
