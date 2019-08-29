package cmd

import (
	"fmt"
	"os"

	"github.com/markbates/pkger"
)

func walk(args []string) error {
	err := pkger.Walk(".", func(path pkger.Path, info os.FileInfo) error {
		fmt.Println(path)
		return nil
	})
	return err
}
