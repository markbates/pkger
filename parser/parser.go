package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
)

var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

func Parse(name string) (Results, error) {
	var r Results
	c, err := pkger.Stat()
	if err != nil {
		return r, err
	}

	if name == "" {
		name = c.ImportPath
	}

	pt, err := pkger.Parse(name)
	if err != nil {
		return r, err
	}
	r.Path = pt

	her, err := pkger.Info(r.Path.Pkg)

	if err != nil {
		return r, err
	}

	m := map[pkger.Path]bool{}

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

	var found []pkger.Path

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

func sourceFiles(pt pkger.Path) ([]pkger.Path, error) {
	var res []pkger.Path

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
		pt := pkger.Path{
			Name: n,
		}
		res = append(res, pt)
		return nil
	})
	return res, err

}

type Results struct {
	Paths []pkger.Path
	Path  pkger.Path
}
