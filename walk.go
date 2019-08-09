package pkger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger/here"
)

type WalkFunc func(Path, os.FileInfo) error

func Walk(p string, wf WalkFunc) error {
	pt, err := Parse(p)
	if err != nil {
		return err
	}
	filesCache.Range(func(k Path, v *File) bool {
		if k.Pkg != pt.Pkg {
			return true
		}
		if err = wf(k, v.info); err != nil {
			if err == filepath.SkipDir {
				return true
			}
			return false
		}
		return true
	})

	if err != nil {
		return err
	}

	var info here.Info
	if pt.Pkg == "." {
		info, err = Stat()
		if err != nil {
			return err
		}
		pt.Pkg = info.ImportPath
	}

	if info.IsZero() {
		info, err = Info(pt.Pkg)
		if err != nil {
			return fmt.Errorf("%s: %s", pt, err)
		}
	}
	fp := filepath.Join(info.Dir, pt.Name)
	err = filepath.Walk(fp, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, info.Dir)
		pt, err := Parse(fmt.Sprintf("%s:%s", pt.Pkg, path))
		if err != nil {
			return err
		}
		return wf(pt, NewFileInfo(fi))
	})

	return err
}
