package parser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/markbates/pkger/paths"
	"github.com/markbates/pkger/pkgs"
)

var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures", "testdata"}

func Parse(name string) (Results, error) {
	var r Results
	if name == "" || name == "." {
		c, err := pkgs.Current()
		if err != nil {
			return r, err
		}
		name = c.ImportPath
	}
	pt, err := paths.Parse(name)
	if err != nil {
		return r, err
	}
	r.Path = pt

	her, err := pkgs.Pkg(r.Path.Pkg)

	if err != nil {
		return r, err
	}

	m := map[paths.Path]bool{}
	root := filepath.Join(her.Dir, r.Path.Name)
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
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".go" {
			return nil
		}

		v, err := NewVisitor(path)
		if err != nil {
			return err
		}

		found, err := v.Run()
		if err != nil {
			return err
		}

		for _, p := range found {
			if _, ok := m[p]; ok {
				continue
			}

			m[p] = true
			if len(p.Name) > 0 {
				continue
			}
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

	var found []paths.Path

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

func sourceFiles(pt paths.Path) ([]paths.Path, error) {
	var res []paths.Path

	her, err := pkgs.Pkg(pt.Pkg)

	if err != nil {
		return res, err
	}
	err = filepath.Walk(her.Dir, func(p string, info os.FileInfo, err error) error {
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

		n := strings.TrimPrefix(strings.TrimPrefix(p, her.Dir), "/")
		pt := paths.Path{
			Name: n,
		}
		res = append(res, pt)
		return nil
	})
	return res, err

}

type Results struct {
	Paths []paths.Path
	Path  paths.Path
}
