package pkger

import (
	"sync"

	"github.com/markbates/pkger/pkging"
)

var current pkging.Pkger
var gil = &sync.RWMutex{}

func Apply(pkg pkging.Pkger, err error) error {
	if err != nil {
		return err
	}
	gil.Lock()
	defer gil.Unlock()
	current = pkging.Wrap(current, pkg)
	return nil
}
