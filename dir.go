package pkger

import (
	"os"
	"path/filepath"
)

func MkdirAll(path string, perm os.FileMode) error {
	pt, err := rootIndex.Parse(path)
	if err != nil {
		return err
	}
	return rootIndex.MkdirAll(pt, perm)
}

func (i *index) MkdirAll(path Path, perm os.FileMode) error {
	root := path.Name

	for root != "" && root != "/" {
		pt := Path{
			Pkg:  path.Pkg,
			Name: root,
		}
		f, err := i.Create(pt)
		if err != nil {
			return err
		}
		f.info.isDir = true
		f.info.mode = perm
		f.info.virtual = true
		if err := f.Close(); err != nil {
			return err
		}
		i.Files.Store(pt, f)
		root = filepath.Dir(root)
	}

	return nil

}
