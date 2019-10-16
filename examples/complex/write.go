package main

import (
	"fmt"

	"github.com/markbates/pkger"
)

// writeFile will write *REAL* files when
// not packaged. When packaged, writeFile
// will write into memory instead.
func writeFile() error {
	// make a folder structure to write into.
	// just like the `os` package directory
	// structures are *NOT* created for you.
	if err := pkger.MkdirAll("/delta/88", 0755); err != nil {
		return err
	}
	// remove the /delta folder and anything
	// underneath it
	defer pkger.RemoveAll("/delta")

	// create a new file under the `/delta/88`
	// directory named `a.car`.
	f, err := pkger.Create("/delta/88/a.car")
	if err != nil {
		return err
	}

	msg := []byte("The Caravan Stops")
	i, err := f.Write(msg)
	if err != nil {
		return err
	}
	if i != len(msg) {
		return fmt.Errorf("expected to write %d bytes, wrote %d instead", len(msg), i)
	}

	// close the file
	if err := f.Close(); err != nil {
		return err
	}

	// stat the new file and get back its os.FileInfo
	info, err := pkger.Stat("/delta/88/a.car")
	if err != nil {
		return err
	}

	fmt.Println("info.Name()\t", info.Name())
	fmt.Println("info.Size()\t", info.Size())
	fmt.Println("info.Mode()\t", info.Mode())
	fmt.Println("info.IsDir()\t", info.IsDir())

	return nil
}
