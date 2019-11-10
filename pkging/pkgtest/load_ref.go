package pkgtest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

func LoadFile(name string, ref *Ref, pkg pkging.Pkger) (os.FileInfo, error) {
	root := filepath.Join(ref.root, name)

	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", name)
	}

	_, err = pkg.Current()
	if err != nil {
		return nil, err
	}

	af, err := os.Open(root)
	if err != nil {
		return nil, err
	}
	defer af.Close()

	bf, err := pkg.Create(name)
	if err != nil {
		return nil, err
	}
	defer bf.Close()

	xp := strings.TrimPrefix(root, filepath.Dir(root))
	xp = filepath.Join(ref.Dir, xp)

	cf, err := os.Create(xp)
	if err != nil {
		return nil, err
	}
	defer cf.Close()

	mw := io.MultiWriter(bf, cf)

	_, err = io.Copy(mw, af)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func LoadFiles(name string, ref *Ref, pkg pkging.Pkger) ([]os.FileInfo, error) {
	var infos []os.FileInfo

	her, err := here.Package("github.com/markbates/pkger")
	if err != nil {
		return nil, err
	}

	root := filepath.Join(ref.root, name)

	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", name)
	}

	her, err = pkg.Current()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		infos = append(infos, info)

		xp := strings.TrimPrefix(path, root)
		xp = filepath.Join(name, xp)

		pt, err := pkg.Parse(xp)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := pkg.MkdirAll(pt.Name, info.Mode()); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Join(her.Dir, xp), info.Mode()); err != nil {
				return err
			}
			return nil
		}

		af, err := os.Open(path)
		if err != nil {
			return err
		}
		defer af.Close()

		bf, err := pkg.Create(pt.Name)
		if err != nil {
			return err
		}
		defer bf.Close()

		xp = filepath.Join(her.Dir, xp)

		cf, err := os.Create(xp)
		if err != nil {
			return err
		}
		defer cf.Close()

		mw := io.MultiWriter(bf, cf)

		_, err = io.Copy(mw, af)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, fmt.Errorf("did not load any infos for %s", name)
	}

	return infos, nil
}
