package here

import (
	"fmt"
	"strings"

	"github.com/markbates/pkger/here/internal/pathparser"
)

func (i Info) Parse(p string) (Path, error) {
	pt := Path{
		Pkg:  i.ImportPath,
		Name: "/",
	}

	res, err := pathparser.Parse(p, []byte(p))
	if err != nil {
		return pt, err
	}

	pp, ok := res.(*pathparser.Path)
	if !ok {
		return pt, fmt.Errorf("expected Path, got %T", res)
	}

	if pp.Pkg != nil {
		pt.Pkg = pp.Pkg.Name
	}

	pt.Name = pp.Name

	her := i
	if pt.Pkg != i.ImportPath {
		her, err = Package(pt.Pkg)
		if err != nil {
			return pt, err
		}
	}

	pt.Name = strings.TrimPrefix(pt.Name, her.Dir)
	pt.Name = strings.ReplaceAll(pt.Name, "\\", "/")

	return pt, nil
}
