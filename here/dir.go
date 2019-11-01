package here

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Dir attempts to gather info for the requested directory.
func Dir(p string) (Info, error) {
	i, err := Cache(p, func(p string) (Info, error) {
		var i Info

		fi, err := os.Stat(p)
		if err != nil {
			return i, err
		}

		if !fi.IsDir() {
			p = filepath.Dir(p)
		}

		pwd, err := os.Getwd()
		if err != nil {
			return i, err
		}

		defer os.Chdir(pwd)

		os.Chdir(p)

		b, err := run("go", "list", "-json")
		if err != nil {

			es := err.Error()
			if !(strings.Contains(es, "cannot find module for path .") || strings.Contains(es, "no Go files")) {
				return i, err
			}

			info, err := Dir(filepath.Dir(p))
			if err != nil {
				return info, err
			}
			i.Module = info.Module

			ph := strings.TrimPrefix(p, info.Module.Dir)

			i.ImportPath = path.Join(info.Module.Path, ph)
			i.Name = path.Base(i.ImportPath)

			ph = filepath.Join(info.Module.Dir, ph)
			i.Dir = ph

			return i, err
		}

		if err := json.Unmarshal(b, &i); err != nil {
			return i, err
		}

		return i, nil
	})

	if err != nil {
		return i, err
	}

	Cache(i.ImportPath, func(p string) (Info, error) {
		return i, nil
	})

	return i, nil
}
