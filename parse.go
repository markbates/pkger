package pkger

import (
	"regexp"
	"strings"
)

var pathrx = regexp.MustCompile("([^:]+)(:(/.+))?")

func Parse(p string) (Path, error) {
	return rootIndex.Parse(p)
}

func build(p, pkg, name string) (Path, error) {
	pt := Path{
		Pkg:  pkg,
		Name: name,
	}

	info, err := Stat()
	if err != nil {
		return pt, err
	}

	if strings.HasPrefix(pt.Pkg, "/") || len(pt.Pkg) == 0 {
		pt.Name = pt.Pkg
		pt.Pkg = info.ImportPath
	}

	if pt.Pkg == pt.Name || len(pt.Name) == 0 {
		pt.Pkg = info.ImportPath
	}

	if !strings.HasPrefix(pt.Name, "/") {
		pt.Name = "/" + pt.Name
	}
	rootIndex.Paths.Store(p, pt)
	return pt, nil
}
