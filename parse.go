package pkger

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var cache = pathsMap{
	data: &sync.Map{},
}

var pathrx = regexp.MustCompile("([^:]+)(:(/.+))?")

func Parse(p string) (Path, error) {
	pt, ok := cache.Load(p)
	if ok {
		return pt, nil
	}
	if len(p) == 0 {
		return build(p, "", "")
	}

	res := pathrx.FindAllStringSubmatch(p, -1)
	if len(res) == 0 {
		return pt, fmt.Errorf("could not parse %q", p)
	}

	matches := res[0]

	if len(matches) != 4 {
		return pt, fmt.Errorf("could not parse %q", p)
	}

	return build(p, matches[1], matches[3])
}

func build(p, pkg, name string) (Path, error) {
	pt := Path{
		Pkg:  pkg,
		Name: name,
	}

	info, err := Current()
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
	cache.Store(p, pt)
	return pt, nil
}
