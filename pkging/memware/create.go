package memware

import (
	"path/filepath"
	"time"

	"github.com/markbates/pkger/pkging"
)

// no such file or directory
func (fx *Warehouse) Create(name string) (pkging.File, error) {
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, err
	}

	her, err := fx.Info(pt.Pkg)
	if err != nil {
		return nil, err
	}

	if _, err := fx.Stat(filepath.Dir(pt.Name)); err != nil {
		return nil, err
	}
	f := &File{
		path: pt,
		her:  her,
		info: &pkging.FileInfo{
			Details: pkging.Details{
				Name:    pt.Name,
				Mode:    0644,
				ModTime: pkging.ModTime(time.Now()),
			},
		},
		pkging: fx,
	}

	fx.files.Store(pt, f)

	return f, nil
}
