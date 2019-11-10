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
	here.ClearCache()
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

	b, err = ioutil.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(dir, "go.mod"))
	if err != nil {
		return nil, err
	}

	if _, err := f.Write(b); err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	// c := exec.Command("cp", "-rv", root, dir)
	// fmt.Println(strings.Join(c.Args, " "))
	// c.Stdout = os.Stdout
	// c.Stderr = os.Stderr
	// c.Stdin = os.Stdin
	// if err := c.Run(); err != nil {
	// 	return nil, err
	// }

	return ref, nil
}
