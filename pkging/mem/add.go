package mem

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/markbates/pkger/pkging"
)

// Add copies the pkging.File into the *Pkger
func (fx *Pkger) Add(files ...pkging.File) error {
	for _, f := range files {
		info, err := f.Stat()
		if err != nil {
			return err
		}

		if f.Path().Pkg == fx.Here.ImportPath {
			dir := filepath.Dir(f.Name())
			if dir != "/" {
				if err := fx.MkdirAll(dir, 0755); err != nil {
					return err
				}
			}
		}

		mf := &File{
			Here:   f.Info(),
			info:   pkging.NewFileInfo(info),
			path:   f.Path(),
			pkging: fx,
		}

		if !info.IsDir() {
			bb := &bytes.Buffer{}
			_, err = io.Copy(bb, f)
			if err != nil {
				return err
			}
			mf.data = bb.Bytes()
		}

		fx.files.Store(mf.Path(), mf)
	}

	return nil
}
