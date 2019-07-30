package main

import (
	"log"
	"os"
)

func main() {
	type ex func([]string) error

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("does not compute")
	}

	var fn ex
	switch args[0] {
	case "walk":
		fn = walk
	case "read":
		fn = read
	case "info":
		fn = info
	case "serve":
		fn = serve
	}
	if fn == nil {
		log.Fatal(args)
	}
	if err := fn(args[1:]); err != nil {
		log.Fatal(err)
	}
}
