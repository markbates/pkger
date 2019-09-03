package mem

import (
	"fmt"

	"github.com/markbates/pkger/pkging"
)

func (fx *Pkger) Open(name string) (pkging.File, error) {
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
		pkging: fx,
		info:   pkging.WithName(f.info.Name(), f.info),
		path:   f.path,
		data:   f.data,
		her:    f.her,
	}

	return nf, nil
}
