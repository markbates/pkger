package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	defer func() {
		c := exec.Command("go", "mod", "tidy", "-v")
		fmt.Println(c.Args)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Run()
	}()

	root, err := New()
	if err != nil {
		log.Fatal(err)
	}

	if err := root.Route(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

// does not computee
