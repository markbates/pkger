package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/pkger"
)

func main() {
	err := pkger.Walk("github.com/gobuffalo/envy", func(path pkger.Path, info os.FileInfo) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
