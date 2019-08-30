package pkger

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Path struct {
	Pkg  string `json:"pkg"`
	Name string `json:"name"`
}

func Parse(p string) (Path, error) {
	p = strings.Replace(p, "\\", "/", -1)
	pt, ok := pathsCache.Load(p)
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

func (p Path) String() string {
	if p.Name == "" {
		p.Name = "/"
	}
	return fmt.Sprintf("%s:%s", p.Pkg, p.Name)
}

func (p Path) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		if st.Flag('+') {
			b, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				return
			}
			fmt.Fprint(st, string(b))
			return
		}
		fmt.Fprint(st, p.String())
	case 'q':
		fmt.Fprintf(st, "%q", p.String())
	default:
		fmt.Fprint(st, p.String())
	}
}

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
