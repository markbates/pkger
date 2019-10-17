package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	return filepath.Walk("./public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Println("Name: ", info.Name())
		fmt.Println("Size: ", info.Size())
		fmt.Println("Mode: ", info.Mode())
		fmt.Println("ModTime: ", info.ModTime())
		fmt.Println()
		return nil
	})

}
