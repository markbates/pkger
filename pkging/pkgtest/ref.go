package pkgtest

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/markbates/pkger/here"
)

type Ref struct {
	here.Info
	root string
}

func NewRef() (*Ref, error) {
	her, err := here.Package("github.com/markbates/pkger")
	if err != nil {
		return nil, err
	}
	root := filepath.Join(
		her.Module.Dir,
		"pkging",
		"pkgtest",
		"testdata",
		"ref",
	)

	return newRef(root)
}

func newRef(root string) (*Ref, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}

	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(root)
	dir = filepath.Join(dir, hex.EncodeToString(b))

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

		root: root,
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	for _, n := range []string{"go.mod", "go.sum"} {
		b, err = ioutil.ReadFile(filepath.Join(root, n))
		if err != nil {
			return nil, err
		}

		f, err := os.Create(filepath.Join(dir, n))
		if err != nil {
			return nil, err
		}

		if _, err := f.Write(b); err != nil {
			return nil, err
		}

		if err := f.Close(); err != nil {
			return nil, err
		}
	}

	return ref, nil
}
