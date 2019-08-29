package cmd

import (
	"fmt"
)

// section: main
func Main(args []string, opts *Options) error {
	if len(args) == 0 {
		return fmt.Errorf("you must pass in an argument")
	}

	opts.Flags.Parse(args)
	args = opts.Flags.Args()

	if opts.Help {
		return Usage(opts.Out(), opts.Flags)
	}

	// s := args[0]
	//
	// u, err := url.Parse(s)
	// if err != nil {
	// 	return err
	// }
	//
	// switch u.Scheme {
	// case "file":
	// 	return File(args, opts)
	// case "http", "https":
	// 	return HTTP(args, opts)
	// }

	return fmt.Errorf("don't know how to handle %s", args)
}

// section: main
