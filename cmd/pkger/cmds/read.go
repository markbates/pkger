package cmds

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/markbates/pkger"
)

type readCmd struct {
	*flag.FlagSet
	JSON bool
}

func (s *readCmd) Name() string {
	return s.Flags().Name()
}

func (r *readCmd) Flags() *flag.FlagSet {
	if r.FlagSet == nil {
		r.FlagSet = flag.NewFlagSet("read", flag.ExitOnError)
		r.FlagSet.BoolVar(&r.JSON, "json", false, "print as JSON")
	}
	return r.FlagSet
}

func (r *readCmd) Exec(args []string) error {
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

		if fi.IsDir() && !r.JSON {
			return fmt.Errorf("can not read a dir %s", a)
		}
		if r.JSON {
			err = json.NewEncoder(os.Stdout).Encode(f)
			if err != nil {
				return err
			}
			continue
		}
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			return err
		}
	}

	return nil
}
