package pkger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	virtual bool
	sys     interface{}
}

func (f *FileInfo) String() string {
	b, _ := json.MarshalIndent(f, "", "  ")
	return string(b)
}

func (f *FileInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":    f.name,
		"size":    f.size,
		"mode":    f.mode,
		"modTime": f.modTime.Format(timeFmt),
		"isDir":   f.isDir,
		"virtual": true,
		"sys":     f.sys,
	})
}

func (f *FileInfo) UnmarshalJSON(b []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	var ok bool

	f.name, ok = m["name"].(string)
	if !ok {
		return fmt.Errorf("could not determine name %q", m["name"])
	}

	size, ok := m["size"].(float64)
	if !ok {
		return fmt.Errorf("could not determine size %q", m["size"])
	}
	f.size = int64(size)

	mode, ok := m["mode"].(float64)
	if !ok {
		return fmt.Errorf("could not determine mode %q", m["mode"])
	}
	f.mode = os.FileMode(mode)

	modTime, ok := m["modTime"].(string)
	if !ok {
		return fmt.Errorf("could not determine modTime %q", m["modTime"])
	}
	t, err := time.Parse(timeFmt, modTime)
	if err != nil {
		return err
	}
	f.modTime = t

	f.isDir, ok = m["isDir"].(bool)
	if !ok {
		return fmt.Errorf("could not determine isDir %q", m["isDir"])
	}
	f.sys = m["sys"]
	f.virtual = true
	return nil
}

func (f *FileInfo) Name() string {
	if !strings.HasPrefix(f.name, "/") {
		f.name = "/" + f.name
	}
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Mode() os.FileMode {
	return f.mode
}

func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *FileInfo) IsDir() bool {
	return f.isDir
}

func (f *FileInfo) Sys() interface{} {
	return f.sys
}

var _ os.FileInfo = &FileInfo{}

func NewFileInfo(info os.FileInfo) *FileInfo {
	fi := &FileInfo{
		name:    info.Name(),
		size:    info.Size(),
		mode:    info.Mode(),
		modTime: info.ModTime(),
		isDir:   info.IsDir(),
		sys:     info.Sys(),
	}
	return fi
}

func WithName(name string, info os.FileInfo) *FileInfo {
	if ft, ok := info.(*FileInfo); ok {
		ft.name = name
		return ft
	}

	fo := NewFileInfo(info)
	fo.name = name
	return fo
}
