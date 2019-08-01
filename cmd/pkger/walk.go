package main

import (
	"fmt"
	"os"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/paths"
)

func walk(args []string) error {
	err := pkger.Walk(".", func(path paths.Path, info os.FileInfo) error {
		fmt.Println(path)
		return nil
	})
	return err
}
