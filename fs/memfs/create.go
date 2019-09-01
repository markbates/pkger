package memfs

import (
	"path/filepath"
	"time"

	"github.com/markbates/pkger/fs"
)

func (fx *FS) Create(name string) (fs.File, error) {
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, err
	}

	her, err := fx.Info(pt.Pkg)
	if err != nil {
		return nil, err
	}
	f := &File{
		path: pt,
		her:  her,
		info: &fs.FileInfo{
			Details: fs.Details{
				Name:    pt.Name,
				Mode:    0644,
				ModTime: fs.ModTime(time.Now()),
			},
		},
		fs: fx,
	}

	fx.files.Store(pt, f)

	dir := filepath.Dir(pt.Name)
	if err := fx.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	return f, nil
}
