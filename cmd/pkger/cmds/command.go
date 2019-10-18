package cmds

import "flag"

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
