package here

import (
	"path/filepath"
	"sync"

	"github.com/markbates/pkger/internal/debug"
)

var curOnce sync.Once
var curErr error
var current Info

func Current() (Info, error) {
	(&curOnce).Do(func() {
		debug.Debug("[HERE] Current")
		b, err := run("go", "env", "GOMOD")
		if err != nil {
			curErr = err
			return
		}
		root := filepath.Dir(string(b))
		i, err := Dir(root)
		if err != nil {
			curErr = err
			return
		}
		current = i
	})

	return current, curErr
}
