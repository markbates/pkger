package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var globalFlags = struct {
	*flag.FlagSet
}{
	FlagSet: flag.NewFlagSet("", flag.ContinueOnError),
}

func main() {

	type ex func([]string) error

	var cmds = map[string]ex{
		"walk":  walk,
		"read":  read,
		"info":  info,
		"serve": serve,
		"list":  list,
		"pack":  pack,
	}
	args := os.Args[1:]

	var fn ex = pack

	if len(args) > 0 {
		var ok bool
		fn, ok = cmds[args[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "couldn't understand args %q\n\n", args)
			fmt.Fprintf(os.Stderr, "the following is a list of available commands:\n\n")
			for k := range cmds {
				fmt.Fprintf(os.Stderr, "\t%s\n", k)
			}
			os.Exit(1)
		}
		args = args[1:]
	}
	if err := fn(args); err != nil {
		log.Fatal(err)
	}
}

// does not computee
