package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

func Parse(her here.Info) (Results, error) {
	var r Results

	name := her.ImportPath

	pt, err := pkger.Parse(name)
	if err != nil {
		return r, err
	}
	r.Path = pt

	m := map[pkging.Path]bool{}

	root := r.Path.Name
	if !strings.HasPrefix(root, string(filepath.Separator)) {
		root = string(filepath.Separator) + root
	}

	if !strings.HasPrefix(root, her.Dir) {
		root = filepath.Join(her.Dir, root)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)

		for _, ig := range DefaultIgnoredFolders {
			if strings.HasPrefix(base, ig) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.IsDir() {
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				her, err = here.Dir(path)
				if err != nil {
					return err
				}
			}
			n := fmt.Sprintf("%s:%s", her.ImportPath, strings.TrimPrefix(path, her.Dir))
			pt, err := pkger.Parse(n)
			if err != nil {
				return err
			}

			m[pt] = true
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".go" {
			return nil
		}

		v, err := newVisitor(path, her)
		if err != nil {
			return err
		}

		found, err := v.Run()
		if err != nil {
			return err
		}

		for _, p := range found {
			if p.Pkg == "." {
				p.Pkg = her.ImportPath
			}
			if _, ok := m[p]; ok {
				continue
			}

			m[p] = true
			found, err := sourceFiles(p)
			if err != nil {
				return err
			}
			for _, pf := range found {
				pf.Pkg = p.Pkg
				m[pf] = true
			}
		}

		return nil
	})

	var found []pkging.Path

	for k := range m {
		if len(k.String()) == 0 {
			continue
		}
		found = append(found, k)
	}
	sort.Slice(found, func(a, b int) bool {
		return found[a].String() <= found[b].String()
	})
	r.Paths = found

	return r, err
}

func sourceFiles(pt pkging.Path) ([]pkging.Path, error) {
	var res []pkging.Path

	her, err := pkger.Info(pt.Pkg)

	if err != nil {
		return res, err
	}

	fp := her.FilePath(pt.Name)
	fi, err := os.Stat(fp)
	if err != nil {
		return res, err
	}
	if !fi.IsDir() {
		return res, nil
	}

	err = filepath.Walk(fp, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(p)

		if base == "." {
			return nil
		}

		for _, ig := range DefaultIgnoredFolders {
			if strings.HasPrefix(base, ig) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		n := strings.TrimPrefix(p, her.Dir)
		n = strings.Replace(n, "\\", "/", -1)
		pt := pkging.Path{
			Name: n,
		}
		res = append(res, pt)
		return nil
	})
	return res, err

}

type Results struct {
	Paths []pkging.Path
	Path  pkging.Path
}
