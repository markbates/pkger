package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {

	defer func() {
		c := exec.Command("go", "mod", "tidy")
		c.Run()
	}()

	type ex func([]string) error

	var cmds = map[string]ex{
		"walk":  walk,
		"read":  read,
		"info":  info,
		"serve": serve,
		"pack":  pack,
	}
	args := os.Args[1:]

	var fn ex = pack

	if len(args) > 0 {
		var ok bool
		fn, ok = cmds[args[0]]
		if ok {
			args = args[1:]
		} else {
			fn = pack
		}
	}
	if err := fn(args); err != nil {
		log.Fatal(err)
	}
}

// does not computee
