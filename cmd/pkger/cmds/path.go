package cmds

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/here"
)

type pathCmd struct {
	*flag.FlagSet
	json bool
	help bool
}

func (s *pathCmd) Name() string {
	return s.Flags().Name()
}

func (c *pathCmd) Flags() *flag.FlagSet {
	if c.FlagSet == nil {
		c.FlagSet = flag.NewFlagSet("path", flag.ExitOnError)
		// c.BoolVar(&c.json, "json", false, "outputs as json")
		c.BoolVar(&c.help, "h", false, "prints help information")
	}
	return c.FlagSet
}

func (c *pathCmd) Exec(args []string) error {

	c.Parse(args)

	if c.help {
		c.Usage()
		return nil
	}

	args = c.Args()
	if len(args) == 0 {
		return fmt.Errorf("you specify at least one path")
	}

	paths := map[string]here.Path{}
	for _, a := range args {
		pt, err := pkger.Parse(a)
		if err != nil {
			return err
		}
		paths[a] = pt
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	return enc.Encode(paths)
}
