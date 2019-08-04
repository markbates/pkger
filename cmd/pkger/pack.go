package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

const outName = "pkged.go"

type packOptions struct {
	*flag.FlagSet
	List bool
}

var packFlags = func() *packOptions {
	rd := &packOptions{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.BoolVar(&rd.List, "list", false, "prints a list of files/dirs to be packaged")
	rd.FlagSet = fs
	return rd
}()

func pack(args []string) error {
	if err := packFlags.Parse(args); err != nil {
		return err
	}
	args = packFlags.Args()

	info, err := pkger.Stat()
	if err != nil {
		return err
	}

	res, err := parser.Parse(info.Dir)
	if err != nil {
		return err
	}

	if packFlags.List {
		fmt.Println(res.Path)

		for _, p := range res.Paths {
			fmt.Printf("  > %s\n", p)
		}
		return nil
	}

	fp := info.FilePath(outName)
	os.RemoveAll(fp)

	if err := Package(fp, res.Paths); err != nil {
		return err
	}

	return nil
}

func Package(out string, paths []pkger.Path) error {
	os.RemoveAll(out)

	f, err := os.Create(out)
	if err != nil {
		return err
	}

	c, err := pkger.Stat()
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "package %s\n\n", c.Name)
	fmt.Fprintf(f, "import \"github.com/markbates/pkger\"\n\n")
	fmt.Fprintf(f, "var _ = pkger.Unpack(`")

	if err := pkger.Pack(f, paths); err != nil {
		return err
	}

	fmt.Fprintf(f, "`)\n")

	return nil
}
