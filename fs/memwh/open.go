package memwh

import (
	"fmt"

	"github.com/markbates/pkger/fs"
)

func (fx *FS) Open(name string) (fs.File, error) {
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, err
	}

	fl, ok := fx.files.Load(pt)
	if !ok {
		return nil, fmt.Errorf("could not open %s", name)
	}
	f, ok := fl.(*File)
	if !ok {
		return nil, fmt.Errorf("could not open %s", name)
	}
	nf := &File{
		fs:   fx,
		info: fs.WithName(f.info.Name(), f.info),
		path: f.path,
		data: f.data,
		her:  f.her,
	}

	return nf, nil
}
