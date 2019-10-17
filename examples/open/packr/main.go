package main

import (
	"log"
	"os"

	"github.com/gobuffalo/packr/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	box := packr.New("demo", "./public")

	b, err := box.Find("index.html")
	if err != nil {
		return err
	}

	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	return nil
}
