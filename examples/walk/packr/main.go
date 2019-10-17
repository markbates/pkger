package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	box := packr.New("demo", "./public")

	return box.Walk(func(path string, f packd.File) error {
		fmt.Println("Name: ", f.Name())
		fmt.Println()
		return nil
	})

}
