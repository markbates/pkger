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
	}
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("does not compute")
	}

	fn, ok := cmds[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "couldn't understand args %q\n\n", args)
		fmt.Fprintf(os.Stderr, "the following is a list of available commands:\n\n")
		for k := range cmds {
			fmt.Fprintf(os.Stderr, "\t%s\n", k)
		}
		os.Exit(1)
	}
	if err := fn(args[1:]); err != nil {
		log.Fatal(err)
	}
}
