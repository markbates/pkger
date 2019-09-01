package memware

import (
	"path/filepath"
	"time"

	"github.com/markbates/pkger/pkging"
)

func (fx *Warehouse) Create(name string) (pkging.File, error) {
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

	dir := filepath.Dir(pt.Name)
	if err := fx.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	return f, nil
}
