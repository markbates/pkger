package pkger

import (
	"regexp"
	"strings"
)

var pathrx = regexp.MustCompile("([^:]+)(:(/.+))?")

func build(p, pkg, name string) (Path, error) {
	pt := Path{
		Pkg:  pkg,
		Name: name,
	}

	current, err := Stat()
	if err != nil {
		return pt, err
	}

	if strings.HasPrefix(pt.Pkg, "/") || len(pt.Pkg) == 0 {
		pt.Name = pt.Pkg
		pt.Pkg = current.ImportPath
	}

	if len(pt.Name) == 0 {
		pt.Name = "/"
	}

	if pt.Pkg == pt.Name {
		pt.Pkg = current.ImportPath
		pt.Name = "/"
	}

	if !strings.HasPrefix(pt.Name, "/") {
		pt.Name = "/" + pt.Name
	}
	pathsCache.Store(p, pt)
	return pt, nil
}
