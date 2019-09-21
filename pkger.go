package pkger

import (
	"os"
	"path/filepath"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/stdos"
)

var _ = func() error {
	return Apply(stdos.New())
}()

func Parse(p string) (pkging.Path, error) {
	return current.Parse(p)
}

func Abs(p string) (string, error) {
	return current.Abs(p)
}

func AbsPath(p pkging.Path) (string, error) {
	return current.AbsPath(p)
}

func Current() (here.Info, error) {
	return current.Current()
}

func Info(p string) (here.Info, error) {
	return current.Info(p)
}

// Create creates the named file with mode 0666 (before umask) - It's actually 0644, truncating it if it already exists. If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR.
func Create(p string) (pkging.File, error) {
	return current.Create(p)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error. The permission bits perm (before umask) are used for all directories that MkdirAll creates. If path is already a directory, MkdirAll does nothing and returns nil.
func MkdirAll(p string, perm os.FileMode) error {
	return current.MkdirAll(p, perm)
}

// Open opens the named file for reading. If successful, methods on the returned file can be used for reading; the associated file descriptor has mode O_RDONLY.
func Open(p string) (pkging.File, error) {
	return current.Open(p)
}

// Stat returns a FileInfo describing the named file.
func Stat(name string) (os.FileInfo, error) {
	return current.Stat(name)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root. All errors that arise visiting files and directories are filtered by walkFn. The files are walked in lexical order, which makes the output deterministic but means that for very large directories Walk can be inefficient. Walk does not follow symbolic links. - That is from the standard library. I know. Their grammar teachers can not be happy with them right now.
func Walk(p string, wf filepath.WalkFunc) error {
	return current.Walk(p, wf)
}

// Remove removes the named file or (empty) directory.
func Remove(name string) error {
	return current.Remove(name)
}

// RemoveAll removes path and any children it contains. It removes everything it can but returns the first error it encounters. If the path does not exist, RemoveAll returns nil (no error).
func RemoveAll(name string) error {
	return current.RemoveAll(name)
}
