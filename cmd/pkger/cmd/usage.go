package cmd

import (
	"flag"
	"fmt"
	"io"
)

func Usage(w io.Writer, f *flag.FlagSet) error {
	fmt.Fprintf(w, "Usage:\n\n")
	fmt.Fprintf(w, "%s [flags] [args...]\n", f.Name())
	f.VisitAll(func(fl *flag.Flag) {
		fmt.Fprintf(w, "\t-%s\t%s (%q)\n", fl.Name, fl.Usage, fl.DefValue)
	})
	return nil
}
