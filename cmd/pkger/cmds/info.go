package cmds

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/markbates/pkger"
)

type infoCmd struct {
	*flag.FlagSet
}

func (s *infoCmd) Name() string {
	return s.Flags().Name()
}

func (r *infoCmd) Flags() *flag.FlagSet {
	if r.FlagSet == nil {
		r.FlagSet = flag.NewFlagSet("info", flag.ExitOnError)
	}
	return r.FlagSet
}

func (r *infoCmd) Exec(args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, a := range args {

		fi, err := pkger.Info(a)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")
		if err := enc.Encode(fi); err != nil {
			return err
		}
	}

	return nil
}
