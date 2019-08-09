package pkger

import (
	"path/filepath"
	"time"
)

func Create(name string) (*File, error) {
	pt, err := Parse(name)
	if err != nil {
		return nil, err
	}

	her, err := Info(pt.Pkg)
	if err != nil {
		return nil, err
	}
	f := &File{
		path: pt,
		her:  her,
		info: &FileInfo{
			name:    pt.Name,
			mode:    0666,
			modTime: time.Now(),
			virtual: true,
		},
	}

	filesCache.Store(pt, f)

	dir := filepath.Dir(pt.Name)
	if err := MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	return f, nil
}
