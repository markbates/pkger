package here

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var cache = &infoMap{
	data: &sync.Map{},
}

func run(n string, args ...string) ([]byte, error) {
	c := exec.Command(n, args...)

	bb := &bytes.Buffer{}
	c.Stdout = bb
	c.Stderr = bb
	err := c.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %q: %s", err, strings.Join(c.Args, " "), bb)
	}

	return bb.Bytes(), nil
}

func Cache(p string, fn func(string) (Info, error)) (Info, error) {
	i, ok := cache.Load(p)
	if ok {
		return i, nil
	}
	i, err := fn(p)
	if err != nil {
		return i, err
	}
	cache.Store(p, i)
	return i, nil
}

func ClearCache() {
	cache = &infoMap{
		data: &sync.Map{},
	}
}

var nonGoDirRx = regexp.MustCompile(`cannot find main|go help modules|go: |build .:|no Go files|can't load package`)
