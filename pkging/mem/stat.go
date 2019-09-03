package mem

import (
	"fmt"
	"os"
)

func (fx *Pkger) Stat(name string) (os.FileInfo, error) {
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, err
	}
	f, ok := fx.files.Load(pt)
	if ok {
		return f.Stat()
	}
	return nil, fmt.Errorf("could not stat %s", pt)
}
