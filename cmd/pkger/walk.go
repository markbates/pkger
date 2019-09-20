package main

import (
	"fmt"
	"os"

	"github.com/markbates/pkger"
)

func walk(args []string) error {
	err := pkger.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	return err
}
