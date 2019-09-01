package memware

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/pkging"
)

func (f *Warehouse) Walk(p string, wf filepath.WalkFunc) error {
	keys := f.files.Keys()

	pt, err := f.Parse(p)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if !strings.HasPrefix(k.Name, pt.Name) {
			continue
		}
		fl, ok := f.files.Load(k)
		if !ok {
			return fmt.Errorf("could not find %s", k)
		}
		fi, err := fl.Stat()
		if err != nil {
			return err
		}

		fi = pkging.WithName(strings.TrimPrefix(k.Name, pt.Name), fi)
		err = wf(k.String(), fi, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
