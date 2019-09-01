package memwh

import (
	"os"
	"path/filepath"
	"time"

	"github.com/markbates/pkger/fs"
)

func (fx *Warehouse) MkdirAll(p string, perm os.FileMode) error {
	path, err := fx.Parse(p)
	if err != nil {
		return err
	}
	root := path.Name

	cur, err := fx.Current()
	if err != nil {
		return err
	}
	for root != "" {
		pt := fs.Path{
			Pkg:  path.Pkg,
			Name: root,
		}
		if _, ok := fx.files.Load(pt); ok {
			root = filepath.Dir(root)
			if root == "/" || root == "\\" {
				break
			}
			continue
		}
		f := &File{
			fs:   fx,
			path: pt,
			her:  cur,
			info: &fs.FileInfo{
				Details: fs.Details{
					Name:    pt.Name,
					Mode:    perm,
					ModTime: fs.ModTime(time.Now()),
				},
			},
		}

		if err != nil {
			return err
		}
		f.info.Details.IsDir = true
		f.info.Details.Mode = perm
		if err := f.Close(); err != nil {
			return err
		}
		fx.files.Store(pt, f)
		root = filepath.Dir(root)
	}

	return nil

}
