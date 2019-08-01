package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/paths"
)

func main() {
	err := pkger.Walk("github.com/gobuffalo/envy", func(path paths.Path, info os.FileInfo) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
