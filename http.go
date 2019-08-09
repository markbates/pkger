package pkger

import (
	"net/http"
	"path"
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
		return nil, err
	}

	fi, err := di.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		di.parent = f.path
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
		info: WithName(f.info.Name(), f.info),
		path: f.path,
		data: f.data,
		her:  f.her,
	}

	return nf, nil
}
