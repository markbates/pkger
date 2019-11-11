package cmds

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
)

type parseCmd struct {
	*flag.FlagSet
	json bool
	help bool
}

func (s *parseCmd) Name() string {
	return s.Flags().Name()
}

func (c *parseCmd) Flags() *flag.FlagSet {
	if c.FlagSet == nil {
		c.FlagSet = flag.NewFlagSet("parse", flag.ExitOnError)
		// c.BoolVar(&c.json, "json", false, "outputs as json")
		c.BoolVar(&c.help, "h", false, "prints help information")
	}
	return c.FlagSet
}

func (c *parseCmd) Exec(args []string) error {

	c.Parse(args)

	if c.help {
		c.Usage()
		return nil
	}

	args = c.Args()
	if len(args) == 0 {
		args = append(args, ".")
	}

	m := map[string]parser.Decls{}

	for _, a := range args {
		var info here.Info
		var err error

		if a == "." {
			info, err = here.Dir(a)
			if err != nil {
				return err
			}
		} else {
			info, err = here.Package(a)
			if err != nil {
				return err
			}

		}
		decls, err := parser.Parse(info)
		if err != nil {
			return err
		}
		m[a] = decls
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	return enc.Encode(m)
}
