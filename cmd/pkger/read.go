package main

import (
	"fmt"
	"io"
	"os"

	"github.com/markbates/pkger"
)

func read(args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, a := range args {
		fmt.Printf("### cmd/pkger/read.go:16 a (%T) -> %q %+v\n", a, a, a)
		f, err := pkger.Open(a)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Println(f.Path())

		fi, err := f.Stat()
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return fmt.Errorf("can not read a dir %s", a)
		}

		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			return err
		}
	}

	return nil
}
