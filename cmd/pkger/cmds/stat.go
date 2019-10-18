package cmds

import (
	"flag"
	"fmt"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

type statCmd struct {
	*flag.FlagSet
}

func (s *statCmd) Name() string {
	return s.Flags().Name()
}

func (f *statCmd) Flags() *flag.FlagSet {
	if f.FlagSet == nil {
		f.FlagSet = flag.NewFlagSet("stat", flag.ExitOnError)
	}
	return f.FlagSet
}

func (f *statCmd) Exec(args []string) error {
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
				fmt.Println(pkging.NewFileInfo(ff))
			}
			continue
		}

		fmt.Println(pkging.NewFileInfo(fi))
	}

	return nil
}
