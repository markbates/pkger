package fs

import (
	"os"
	"path/filepath"

	"github.com/markbates/pkger/here"
)

type FileSystem interface {
	Create(name string) (File, error)
	Current() (here.Info, error)
	Info(p string) (here.Info, error)
	MkdirAll(p string, perm os.FileMode) error
	Open(name string) (File, error)
	Parse(p string) (Path, error)
	ReadFile(s string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
	Walk(p string, wf filepath.WalkFunc) error
}
