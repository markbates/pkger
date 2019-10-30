package costello

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/pkging"
)

func LoadRef(ref *Ref, pkg pkging.Pkger) error {
	return filepath.Walk(ref.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		of, err := os.Open(path)
		if err != nil {
			return err
		}
		defer of.Close()

		if a, ok := pkg.(pkging.Adder); ok {
			return a.Add(of)
		}

		path = strings.TrimPrefix(path, ref.Dir)

		pt, err := pkg.Parse(path)
		if err != nil {
			return err
		}

		if err := pkg.MkdirAll(filepath.Dir(pt.Name), 0755); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		f, err := pkg.Create(pt.String())
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, of); err != nil {
			return err
		}
		return nil
	})
}
