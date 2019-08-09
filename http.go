package pkger

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (f *File) Open(name string) (http.File, error) {
	pt, err := Parse(name)
	if err != nil {
		return nil, err
	}

	if pt == f.path {
		return f, nil
	}

	pt.Name = path.Join(f.Path().Name, pt.Name)

	di, err := Open(pt.String())
	if err != nil {
		if filepath.Base(name) == "index.html" {
			if _, ok := err.(*os.PathError); ok {
				return f, nil
			}
		}
		return f, err
	}
	return di, nil
}

func Open(name string) (*File, error) {
	pt, err := Parse(name)
	if err != nil {
		return nil, err
	}

	dubeg("Open", pt.String())
	f, ok := filesCache.Load(pt)
	if !ok {
		return openDisk(pt)
	}
	nf := &File{
		info: WithName(strings.TrimPrefix(f.info.Name(), "/"), f.info),
		path: f.path,
		data: f.data,
		her:  f.her,
	}

	return nf, nil
}
