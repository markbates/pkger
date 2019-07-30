package main

import (
	"fmt"

	"github.com/markbates/pkger"
)

func info(args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, a := range args {
		f, err := pkger.Open(a)
		if err != nil {
			return err
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return err
		}

		if fi.IsDir() {
			files, err := f.Readdir(-1)
			if err != nil {
				return err
			}
			for _, ff := range files {
				fmt.Println(pkger.NewFileInfo(ff))
			}
			continue
		}

		fmt.Println(pkger.NewFileInfo(fi))
	}

	return nil
}
