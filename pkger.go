package pkger

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

func modRoot() (string, error) {
	c := exec.Command("go", "env", "GOMOD")
	b, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return "", fmt.Errorf("the `go env GOMOD` was empty/modules are required")
	}

	return filepath.Dir(string(b)), nil
}

func Getwd() (string, error) {
	return modRoot()
}

func Open(p string) (*File, error) {
	pt, err := Parse(p)
	if err != nil {
		return nil, err
	}
	return rootIndex.Open(pt)
}
