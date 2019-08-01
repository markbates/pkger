package main

import (
	"fmt"

	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkgs"
)

func list(args []string) error {
	info, err := pkgs.Current()
	if err != nil {
		return err
	}
	res, err := parser.Parse(info.Dir)
	if err != nil {
		return err
	}

	fmt.Println(res.Path)

	for _, p := range res.Paths {
		fmt.Printf("  > %s\n", p)
	}
	return nil
}
