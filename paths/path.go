package paths

import (
	"fmt"
	"strings"
)

type Path struct {
	Pkg  string `json:"pkg"`
	Name string `json:"name"`
}

func (p Path) String() string {
	if len(p.Pkg) == 0 {
		return p.Name
	}
	if len(p.Name) == 0 {
		return p.Pkg
	}
	return fmt.Sprintf("%s:/%s", p.Pkg, p.Name)
}

func Parse(p string) (Path, error) {
	var pt Path
	res := strings.Split(p, ":")

	if len(res) < 1 {
		return pt, fmt.Errorf("could not parse %q (%d)", res, len(res))
	}
	if len(res) == 1 {
		if strings.HasPrefix(res[0], "/") {
			pt.Name = res[0]
		} else {
			pt.Pkg = res[0]
		}
	} else {
		pt.Pkg = res[0]
		pt.Name = res[1]
	}
	pt.Name = strings.TrimPrefix(pt.Name, "/")
	pt.Pkg = strings.TrimPrefix(pt.Pkg, "/")
	return pt, nil
}
