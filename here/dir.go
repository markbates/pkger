package here

import (
	"encoding/json"
	"fmt"
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
			if !(strings.Contains(es, "cannot find module for path .") || strings.Contains(es, "no Go files") || strings.Contains(es, "can't load package")) {
				return i, err
			}

			if strings.Contains(es, "can't load package: package .") {
				if _, err := os.Stat(fmt.Sprintf("%s/go.mod", p)); err == nil {
					var mod Module
					bm, err := run ("go", "list", "-m", "-json")
					if err != nil {
						return i, err
					}

					if err := json.Unmarshal(bm, &mod); err != nil {
						return i, err
					}
					info := NewInfoFromPath(p, mod)
					prepareInfo(p, info, &i)

					return i, err
				}
			}

			info, err := Dir(filepath.Dir(p))
			if err != nil {
				return info, err
			}

			prepareInfo(p, info, &i)
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

func prepareInfo(p string, info Info, target *Info) {
	target.Module = info.Module
	ph := strings.TrimPrefix(p, target.Module.Dir)

	target.ImportPath = path.Join(info.Module.Path, ph)
	target.Name = path.Base(target.ImportPath)

	ph = filepath.Join(info.Module.Dir, ph)
	target.Dir = ph
}