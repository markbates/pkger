package pkger

import (
	"os"
	"path/filepath"
	"time"
)

func MkdirAll(p string, perm os.FileMode) error {
	path, err := Parse(p)
	if err != nil {
		return err
	}
	root := path.Name

	cur, err := Stat()
	if err != nil {
		return err
	}
	for root != "" && root != "/" {
		pt := Path{
			Pkg:  path.Pkg,
			Name: root,
		}
		if _, ok := filesCache.Load(pt); ok {
			root = filepath.Dir(root)
			continue
		}
		f := &File{
			path: pt,
			her:  cur,
			info: &FileInfo{
				name:    pt.Name,
				mode:    0666,
				modTime: time.Now(),
				virtual: true,
			},
		}

		filesCache.Store(pt, f)
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
