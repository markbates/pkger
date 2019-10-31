package costello

import (
	"os"
	"path/filepath"
	"runtime"

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
			Name:       "app",
			Module: here.Module{
				Main:      true,
				Path:      "app",
				Dir:       dir,
				GoMod:     filepath.Join(dir, "go.mod"),
				GoVersion: runtime.Version(),
			},
		},
	}

	return ref, nil
}
