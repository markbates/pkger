package main

import (
	"flag"
	"fmt"

	"github.com/markbates/pkger"
)

type infoCmd struct {
	*flag.FlagSet
}

func (s *infoCmd) Name() string {
	return s.Flags().Name()
}

func (f *infoCmd) Flags() *flag.FlagSet {
	if f.FlagSet == nil {
		f.FlagSet = flag.NewFlagSet("info", flag.ExitOnError)
	}
	return f.FlagSet
}

func (f *infoCmd) Exec(args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, a := range args {
		f, err := pkger.Open(a)
		if err != nil {
			return err
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return err
		}

		if fi.IsDir() {
			files, err := f.Readdir(-1)
			if err != nil {
				return err
			}
			for _, ff := range files {
				fmt.Println(pkger.NewFileInfo(ff))
			}
			continue
		}

		fmt.Println(pkger.NewFileInfo(fi))
	}

	return nil
}
