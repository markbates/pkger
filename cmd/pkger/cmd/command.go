package cmd

import (
	"flag"
	"fmt"
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
