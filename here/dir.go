package here

import (
	"encoding/json"
	"fmt"
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
		// go: cannot find main module; see 'go help modules'
		// build .: cannot find module for path .
		// no Go files in
		if err != nil {
			return fromNonGoDir(p)
		}

		if err := json.Unmarshal(b, &i); err != nil {
			return i, err
		}

		return i, nil
	})

	if err != nil {
		return i, err
	}

	return Cache(i.ImportPath, func(p string) (Info, error) {
		return i, nil
	})

}

func dir(p string) (Info, error) {
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
	// go: cannot find main module; see 'go help modules'
	// build .: cannot find module for path .
	// no Go files in
	if err != nil {
		return fromNonGoDir(p)
	}

	if err := json.Unmarshal(b, &i); err != nil {
		return i, err
	}

	return i, nil
}

func fromNonGoDir(dir string) (Info, error) {
	fmt.Println(">>>TODO here/dir.go:59: dir ", dir)
	i := Info{
		Dir: dir,
	}

	b, err := run("go", "list", "-json", "-m")
	if err != nil {
		return i, err
	}

	if err := json.Unmarshal(b, &i.Module); err != nil {
		return i, err
	}

	return i, err
}
