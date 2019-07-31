package pkger

import (
	"fmt"
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
	return rootIndex.Parse(p)
}
