package here

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/markbates/pkger/internal/takeon/github.com/markbates/hepa"
	"github.com/markbates/pkger/internal/takeon/github.com/markbates/hepa/filters"
)

// Info represents details about the directory/package
type Info struct {
	Dir        string
	ImportPath string
	Name       string
	Imports    []string
	Module     Module
}

func (fi Info) MarshalJSON() ([]byte, error) {
	mm := map[string]interface{}{
		"ImportPath": fi.ImportPath,
		"Name":       fi.Name,
		"Imports":    fi.Imports,
		"Module":     fi.Module,
	}

	hep := hepa.New()
	hep = hepa.With(hep, filters.Home())
	hep = hepa.With(hep, filters.Golang())

	cm := map[string]string{
		"Dir": fi.Dir,
	}

	for k, v := range cm {
		b, err := hep.Filter([]byte(v))
		if err != nil {
			return nil, err
		}
		mm[k] = string(b)
	}

	return json.Marshal(mm)
}

func (i Info) FilePath(paths ...string) string {
	res := []string{i.Dir}
	for _, p := range paths {
		p = strings.TrimPrefix(p, i.Dir)
		p = strings.TrimPrefix(p, "/")
		if runtime.GOOS == "windows" {
			p = strings.Replace(p, "/", "\\", -1)
		}
		res = append(res, p)
	}
	return filepath.Join(res...)
}

func (i Info) Open(p string) (*os.File, error) {
	return os.Open(i.FilePath(p))
}

// ModuleName returns the name of the current
// module, or if not using modules, the current
// package. These *might* not match.
func (i Info) ModuleName() string {
	if i.Mods() {
		return i.Module.Path
	}
	return i.ImportPath
}

// IsZero checks if the type has been filled
// with rich chocolately data goodness
func (i Info) IsZero() bool {
	return i.String() == Info{}.String()
}

// Mods returns whether Go modules are used
// in this directory/package.
func (i Info) Mods() bool {
	return !i.Module.IsZero()
}

func (i Info) String() string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	s := string(b)
	return s
}

func (i Info) Parse(p string) (Path, error) {
	p = strings.TrimSpace(p)
	p = filepath.Clean(p)
	p = strings.TrimPrefix(p, i.Dir)

	p = strings.Replace(p, "\\", "/", -1)
	p = strings.TrimSpace(p)

	if len(p) == 0 || p == ":" {
		return i.build("", "", "")
	}

	res := pathrx.FindAllStringSubmatch(p, -1)
	if len(res) == 0 {
		return Path{}, fmt.Errorf("could not parse %q", p)
	}

	matches := res[0]

	if len(matches) != 4 {
		return Path{}, fmt.Errorf("could not parse %q", p)
	}

	return i.build(p, matches[1], matches[3])
}

func (i Info) build(p, pkg, name string) (Path, error) {
	pt := Path{
		Pkg:  pkg,
		Name: name,
	}

	if strings.HasPrefix(pt.Pkg, "/") || len(pt.Pkg) == 0 {
		pt.Name = pt.Pkg
		pt.Pkg = i.ImportPath
	}

	if len(pt.Name) == 0 {
		pt.Name = "/"
	}

	if pt.Pkg == pt.Name {
		pt.Pkg = i.ImportPath
		pt.Name = "/"
	}

	if !strings.HasPrefix(pt.Name, "/") {
		pt.Name = "/" + pt.Name
	}
	pt.Name = strings.TrimPrefix(pt.Name, i.Dir)
	return pt, nil
}

var pathrx = regexp.MustCompile("([^:]+)(:(/.+))?")
