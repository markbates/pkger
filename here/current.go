package here

import (
	"path/filepath"
	"sync"
)

var curOnce sync.Once
var curErr error
var current Info

func Current() (Info, error) {
	(&curOnce).Do(func() {
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
