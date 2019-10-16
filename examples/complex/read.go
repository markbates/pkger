package main

import (
	"io"
	"os"

	"github.com/markbates/pkger"
)

func readFile() error {
	f, err := pkger.Open("/go.mod")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(os.Stdout, f)
	return err
}
