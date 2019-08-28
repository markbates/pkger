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
		return nil, err
	}

	fi, err := di.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		di.parent = f.path
		di.excludes = f.excludes
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

func openDisk(pt Path) (*File, error) {
	dubeg("openDisk", pt.String())
	info, err := Info(pt.Pkg)
	if err != nil {
		return nil, err
	}
	fp := info.Dir
	if len(pt.Name) > 0 {
		fp = filepath.Join(fp, pt.Name)
	}

	fi, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}
	f := &File{
		info: WithName(strings.TrimPrefix(pt.Name, "/"), NewFileInfo(fi)),
		her:  info,
		path: pt,
	}
	return f, nil
}
