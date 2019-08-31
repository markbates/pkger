package fstest

import (
	"bytes"
	"io"

	"github.com/markbates/pkger/fs"
)

type TestFile struct {
	Name string
	Data []byte
}

func (t TestFile) Create(fx fs.FileSystem) error {
	f, err := fx.Create(t.Name)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(t.Data))
	if err != nil {
		return err
	}
	return f.Close()
}

type TestFiles map[string]TestFile

func (t TestFiles) Create(fx fs.FileSystem) error {
	for k, f := range t {
		f.Name = k
		if err := f.Create(fx); err != nil {
			return err
		}
	}
	return nil
}
