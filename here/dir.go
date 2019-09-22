package here

import (
	"encoding/json"
	"os"
	"path/filepath"
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
