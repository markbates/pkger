package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/markbates/pkger"
)

type readOptions struct {
	*flag.FlagSet
	JSON bool
}

var readFlags = func() *readOptions {
	rd := &readOptions{}
	fs := flag.NewFlagSet("read", flag.ExitOnError)
	fs.BoolVar(&rd.JSON, "json", false, "print as JSON")
	rd.FlagSet = fs
	return rd
}()

func read(args []string) error {
	if err := readFlags.Parse(args); err != nil {
		return err
	}
	args = readFlags.Args()

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

		if fi.IsDir() && !readFlags.JSON {
			return fmt.Errorf("can not read a dir %s", a)
		}
		if readFlags.JSON {
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
