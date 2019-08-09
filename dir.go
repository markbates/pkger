package pkger

import (
	"os"
	"path/filepath"
)

func MkdirAll(p string, perm os.FileMode) error {
	path, err := Parse(p)
	if err != nil {
		return err
	}
	root := path.Name

	for root != "" && root != "/" {
		pt := Path{
			Pkg:  path.Pkg,
			Name: root,
		}
		f, err := Create(pt.String())
		if err != nil {
			return err
		}
		f.info.isDir = true
		f.info.mode = perm
		f.info.virtual = true
		if err := f.Close(); err != nil {
			return err
		}
		filesCache.Store(pt, f)
		root = filepath.Dir(root)
	}

	return nil

}
