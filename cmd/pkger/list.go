package main

import (
	"fmt"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

func list(args []string) error {
	info, err := pkger.Current()
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
