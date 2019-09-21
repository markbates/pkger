package pkger

import (
	"log"
	"os"
	"sync"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgutil"
)

var current pkging.Pkger
var gil = &sync.RWMutex{}

func Apply(pkg pkging.Pkger, err error) error {
	if err := pkgutil.Dump(os.Stdout, pkg); err != nil {
		log.Fatal(err)
		return err
	}
	gil.Lock()
	defer gil.Unlock()
	current = pkging.Wrap(current, pkg)
	return nil
}
