package pkger

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	di, err := rootIndex.Open(pt)
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

func (i *index) Open(pt Path) (*File, error) {
	i.debug("Open", pt.String())
	f, ok := i.Files.Load(pt)
	if !ok {
		return i.openDisk(pt)
	}
	nf := &File{
		info: f.info,
		path: f.path,
		data: f.data,
		her:  f.her,
	}

	return nf, nil
}
