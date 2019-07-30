package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/pkger"
)

func main() {
	err := pkger.Walk("github.com/gobuffalo/buffalo", func(path pkger.Path, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
