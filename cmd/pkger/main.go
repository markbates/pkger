package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/markbates/pkger/cmd/pkger/cmds"
)

func main() {
	clean := func() {
		c := exec.Command("go", "mod", "tidy")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Run()
	}
	defer clean()

	defer func() {
		if err := recover(); err != nil {
			clean()
			log.Fatal(err)
		}
	}()

	if err := run(); err != nil {
		clean()
		log.Fatal(err)
	}
}

func run() error {
	root, err := cmds.New()
	if err != nil {
		return err
	}

	return root.Route(os.Args[1:])
}

// does not compute
