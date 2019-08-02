package pkger

import (
	"os"
)

type WalkFunc func(Path, os.FileInfo) error

func Walk(p string, wf WalkFunc) error {
	pt, err := Parse(p)
	if err != nil {
		return err
	}
	return rootIndex.Walk(pt, wf)
}
