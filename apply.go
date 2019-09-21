package pkger

import (
	"os"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgutil"
)

var current pkging.Pkger

func Apply(pkg pkging.Pkger, err error) error {
	if err := pkgutil.Dump(os.Stdout, pkg); err != nil {
		return err
	}
	current = pkging.Wrap(current, pkg)
	return nil
}
