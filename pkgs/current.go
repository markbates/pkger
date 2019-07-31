package pkgs

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/here"
)

func Pkg(p string) (here.Info, error) {
	return here.Cache(p, here.Package)
}

func Dir(p string) (here.Info, error) {
	return here.Cache(p, here.Dir)
}

func Current() (here.Info, error) {
	return Dir(".")
}

func Stat(info here.Info, p string) (os.FileInfo, error) {
	return os.Stat(filepath.Join(info.Dir, p))
}
