package cmd

import "flag"

type Options struct {
	IO
	Flags *flag.FlagSet
	Help  bool
}

func NewOptions(gio IO) *Options {
	g := &Options{
		IO:    gio,
		Flags: flag.NewFlagSet("pkger", flag.ExitOnError),
	}
	g.Flags.BoolVar(&g.Help, "h", false, "print help")
	return g
}
