package here

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Info represents details about the directory/package
type Info struct {
	Dir         string
	ImportPath  string
	Name        string
	Doc         string
	Target      string
	Root        string
	Match       []string
	Stale       bool
	StaleReason string
	GoFiles     []string
	Imports     []string
	Deps        []string
	TestGoFiles []string
	TestImports []string
	Module      Module
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
