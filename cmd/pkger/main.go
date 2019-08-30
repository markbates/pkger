package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type command interface {
	Name() string
	execer
	flagger
}

type execer interface {
	Exec([]string) error
}

type flagger interface {
	Flags() *flag.FlagSet
}
type arrayFlags []string

func (i arrayFlags) String() string {
	return fmt.Sprintf("%s", []string(i))
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {

	defer func() {
		c := exec.Command("go", "mod", "tidy")
		c.Run()
	}()

	root := &packCmd{}
	cmds := []command{
		root, &readCmd{}, &serveCmd{}, &infoCmd{},
	}
	sort.Slice(cmds, func(a, b int) bool {
		return cmds[a].Name() <= cmds[b].Name()
	})
	root.Flags().Usage = func() {
		for _, c := range cmds {
			fg := c.Flags()
			fmt.Fprintf(os.Stderr, "%s:\n", fg.Name())
			fg.PrintDefaults()
		}
	}

	root.Parse(os.Args[1:])
	args := root.Args()
	var ex command = root

	if len(args) > 0 {
		k := args[0]
		for _, c := range cmds {
			if k == strings.TrimPrefix(c.Name(), "pkger ") {
				ex = c
				args = args[1:]
				break
			}
		}
	}

	flg := ex.Flags()
	flg.Parse(args)
	args = flg.Args()

	if err := ex.Exec(args); err != nil {
		log.Fatal(err)
	}
}

// does not computee
