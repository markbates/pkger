package main

import (
	"log"
	"os"

	"app/actions"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := actions.WalkTemplates(os.Stdout); err != nil {
		return err
	}
	return nil
}