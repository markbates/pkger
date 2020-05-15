package pkgutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

func Dump(w io.Writer, pkg pkging.Pkger) error {
	d := struct {
		Type  string
		Info  here.Info
		Paths map[string]os.FileInfo
	}{
		Type:  fmt.Sprintf("%T", pkg),
		Paths: map[string]os.FileInfo{},
	}

	info, err := pkg.Current()
	if err != nil {
		return err
	}
	d.Info = info

	err = pkg.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		d.Paths[path] = info
		return nil
	})
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")

	if err := enc.Encode(d); err != nil {
		return err
	}

	return nil
}
