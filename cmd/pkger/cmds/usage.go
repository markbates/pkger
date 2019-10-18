package cmds

import (
	"flag"
	"fmt"
	"io"
)

func Usage(w io.Writer, f *flag.FlagSet) func() {
	return func() {
		pre := "pkger "
		if f.Name() == "pkger" {
			pre = ""
		}
		fmt.Fprintf(w, "%s%s [flags] [args...]\n", pre, f.Name())
		f.VisitAll(func(fl *flag.Flag) {
			fmt.Fprintf(w, "\t-%s\t%s (%q)\n", fl.Name, fl.Usage, fl.DefValue)
		})
		fmt.Fprintln(w)
	}
}
