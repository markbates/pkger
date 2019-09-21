package pkging

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type Details struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime ModTime     `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
	Sys     interface{} `json:"sys"`
}
type FileInfo struct {
	Details `json:"details"`
}

func (f *FileInfo) String() string {
	b, _ := json.MarshalIndent(f, "", "  ")
	return string(b)
}

func (f *FileInfo) Name() string {
	return f.Details.Name
}

func (f *FileInfo) Size() int64 {
	return f.Details.Size
}

func (f *FileInfo) Mode() os.FileMode {
	return f.Details.Mode
}

func (f *FileInfo) ModTime() time.Time {
	return time.Time(f.Details.ModTime)
}

func (f *FileInfo) IsDir() bool {
	return f.Details.IsDir
}

func (f *FileInfo) Sys() interface{} {
	return f.Details.Sys
}

var _ os.FileInfo = &FileInfo{}

func NewFileInfo(info os.FileInfo) *FileInfo {
	fi := &FileInfo{
		Details: Details{
			Name:    cleanName(info.Name()),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: ModTime(info.ModTime()),
			IsDir:   info.IsDir(),
			Sys:     info.Sys(),
		},
	}
	return fi
}

func WithName(name string, info os.FileInfo) *FileInfo {
	s := cleanName(name)

	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	fo := NewFileInfo(info)
	fo.Details.Name = cleanName(name)
	return fo
}

func WithRelName(name string, info os.FileInfo) *FileInfo {
	s := cleanName(name)
	s = strings.TrimPrefix(s, "/")

	fo := NewFileInfo(info)
	fo.Details.Name = cleanName(s)
	return fo
}

func cleanName(s string) string {
	if strings.Contains(s, "\\") {
		s = strings.Replace(s, "\\", "/", -1)
	}
	return s
}
