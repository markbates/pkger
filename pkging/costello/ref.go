package costello

import (
	"os"
	"path/filepath"

	"github.com/markbates/pkger/here"
)

type Ref struct {
	here.Info
}

func NewRef() (*Ref, error) {

	her, err := here.Package("github.com/markbates/pkger")
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(
		her.Module.Dir,
		"pkging",
		"costello",
		"testdata",
		"ref")

	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	ref := &Ref{
		Info: here.Info{
			ImportPath: "app",
			Dir:        dir,
			Module: here.Module{
				Path: "app",
				Dir:  dir,
			},
		},
	}

	return ref, nil
}
