package cmd

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/markbates/pkger"
)

type serveCmd struct {
	*flag.FlagSet
	excludes arrayFlags
}

func (s *serveCmd) Name() string {
	return s.Flags().Name()
}

func (f *serveCmd) Flags() *flag.FlagSet {
	if f.FlagSet == nil {
		f.FlagSet = flag.NewFlagSet("serve", flag.ExitOnError)
		f.Var(&f.excludes, "exclude", "slice of regexp patterns to exclude")
	}
	return f.FlagSet
}

var defaultExcludes = []string{"testdata", "node_modules", "(\\/|\\\\)_.+", "(\\/|\\\\)\\.git.*", ".DS_Store"}

func (s *serveCmd) Exec(args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	f, err := pkger.Open(args[0])
	if err != nil {
		return err
	}

	ex := append(defaultExcludes, s.excludes...)
	if err := pkger.Exclude(f, ex...); err != nil {
		return err
	}
	defer f.Close()
	fmt.Println(f.Path())

	return http.ListenAndServe(":3000", http.FileServer(f))
}
