package stdos

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/costello"
)

func NewTemp(ref *costello.Ref) (pkging.Pkger, error) {
	dir, err := ioutil.TempDir("", "stdos")
	if err != nil {
		return nil, err
	}

	info := here.Info{
		Module:     ref.Module,
		ImportPath: ref.ImportPath,
		Name:       ref.Name,
		Dir:        dir,
	}
	info.Module.Dir = dir
	info.Module.GoMod = filepath.Join(dir, "go.mod")
	return New(info)
}

func Test_Pkger(t *testing.T) {
	costello.All(t, func(ref *costello.Ref) (pkging.Pkger, error) {
		return NewTemp(ref)
	})
}
