// +build debug

package pkger

import (
	"os"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgutil"
)

func Apply(pkg pkging.Pkger, err error) error {
	gil.Lock()
	defer gil.Unlock()
	if err != nil {
		return err
	}
	if err := pkgutil.Dump(os.Stdout, pkg); err != nil {
		return err
	}
	current = pkging.Wrap(current, pkg)
	return nil
}
