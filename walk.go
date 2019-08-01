package pkger

import (
	"os"

	"github.com/markbates/pkger/paths"
)

type WalkFunc func(paths.Path, os.FileInfo) error

func Walk(p string, wf WalkFunc) error {
	pt, err := paths.Parse(p)
	if err != nil {
		return err
	}
	return rootIndex.Walk(pt, wf)
}
